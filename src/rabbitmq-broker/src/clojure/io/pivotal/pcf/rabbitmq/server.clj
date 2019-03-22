(ns io.pivotal.pcf.rabbitmq.server
  (:require [taoensso.timbre :as log]
            [io.pivotal.pcf.rabbitmq.config :as cfg]
            [io.pivotal.pcf.rabbitmq.resources :as rs]
            [clojure.java.io :as io]
            beckon
            [ring.adapter.jetty9 :refer [run-jetty]]
            [compojure.core :refer [defroutes GET PUT DELETE]]
            [compojure.route :as rt]
            [ring.util.response :refer [response status]]
            [cheshire.core :as json]
            [ring.middleware.json :refer [wrap-json-response]]
            [ring.middleware.basic-authentication :refer [wrap-basic-authentication]]
            [clojure.string :as string]
            [clj-http.client :as httpc]
            [langohr.http :as hc])
  (:use [slingshot.slingshot :only [throw+ try+]])
  (:import java.io.File
           java.lang.management.ManagementFactory)
  )

;;
;; Implementation
;;

(defn ^{:private true} log-exception
  [^Exception e]
  (log/errorf "Caught an exception during boot: %s (%s)" (.getMessage e) (.getClass e)) e)

(defn initialize-logger
  [config]
  (log/set-level! (keyword (cfg/log-level config))))

(defn announce-start
  [config]
  (log/infof "Starting. CC endpoint: %s"
             (cfg/cc-endpoint config)))

(defn log-if-using-tls
  [config]
  (if (cfg/using-tls? config)
    (log/infof "Will use HTTPS to talk to RabbitMQ HTTP API as %s..." (cfg/rabbitmq-administrator))
    (log/infof "Will use HTTP (not HTTPS) to talk to RabbitMQ HTTP API as %s..." (cfg/rabbitmq-administrator))))

(declare shutdown)
(defn install-signal-traps
  []
  (let [xs #{shutdown}]
    (reset! (beckon/signal-atom "INT")  xs)
    (reset! (beckon/signal-atom "TERM") xs)))

(defn init-rabbitmq-connection!
  [config]
  (let [uri    (get (cfg/rabbitmq-administrator-uris config) 0)
        uname  (cfg/rabbitmq-administrator config)
        pwd    (cfg/rabbitmq-administrator-password config)
        opts   (if (cfg/using-tls?)
                 ;; don't perform peer verification with RabbitMQ nodes
                 {:insecure? true}
                 {})]
    (hc/connect! uri uname pwd opts)))

(defn wrap-request-logging
  [f]
  (fn [{:keys [request-method uri] :as req}]
    (let [start (System/currentTimeMillis)
          res   (f req)
          end   (System/currentTimeMillis)
          t     (- end start)]
      (log/infof "%s %s %d %d (in %d ms)"
                 (.toUpperCase ^String (name request-method))
                 uri
                 (:status res)
                 (count (:body res))
                 t)
      res)))

(defmacro defresponder
  "Defines a response helper function that has 2 arities:

   * 0-arity responds with an empty body
   * 1-arity responds with the argument as body"
  [name status]
  `(defn ~name
     ([]
        (~name {}))
     ([body#]
        (-> (response body#)
            (status ~status)))))

(defresponder ok             200)
(defresponder created        201)
(defresponder bad-request    400)
(defresponder conflict       409)
(defresponder gone           410)
(defresponder internal-error 500)

;;
;; Routes
;;

(defn forward-request-put
  [req]
  (let [headers (select-keys (get req :headers) ["authorization"])
        endpoint (get req :uri)
        body (slurp (get req :body))]
    (try+
      (httpc/put (format "http://localhost:8901%s" endpoint) {:body body :headers (assoc headers :X-Broker-API-Version "2.14")})
    (catch Object e
      (log/infof "forward-request failed for: %s, %s, %s" endpoint headers body)
      e
    ))))

(defn forward-request-get
  [req]
  (let [headers (select-keys (get req :headers) ["authorization"])
        endpoint (get req :uri)]
    (try+
      (httpc/get (format "http://localhost:8901%s" endpoint) {:headers (assoc headers :X-Broker-API-Version "2.14")})
    (catch Object e
      (log/infof "forward-request failed for: %s, %s" endpoint headers)
      e
    ))))

(defn forward-request-delete
  [req]
  (let [headers (select-keys (get req :headers) ["authorization"])
        endpoint (get req :uri),
        query-string (get req :query-string)]
    (try+
      (httpc/delete (format "http://localhost:8901%s?%s" endpoint query-string) {:headers (assoc headers :X-Broker-API-Version "2.14")})
    (catch Object e
      (log/infof "forward-request failed for: %s, %s" endpoint headers)
      e
    ))))

(defn show-raw-config
  [_]
  (let [pretty-printed (json/generate-string (cfg/serializable-config) {:pretty true})]
    (ok pretty-printed)))

(defn show-cf-api-version
  [_]
  (ok "2.0"))

(defroutes broker-v2-routes
  (GET    "/v2/catalog"               req forward-request-get)
  (PUT    "/v2/service_instances/:id" req forward-request-put)
  (DELETE "/v2/service_instances/:id" req forward-request-delete)
  (PUT    "/v2/service_instances/:instance_id/service_bindings/:id" req forward-request-put)
  (DELETE "/v2/service_instances/:instance_id/service_bindings/:id" req forward-request-delete)
  (GET    "/ops/config"               req show-raw-config)
  (GET    "/ops/cf/api/version"       req show-cf-api-version))

(defn start-http-server
  [config]
  (run-jetty (-> broker-v2-routes
                 wrap-json-response
                 (wrap-basic-authentication cfg/authenticated?)
                 wrap-request-logging)
             {:port 4567
              :join? false}))

;;
;; API
;;

(defn start
  [config]
  (initialize-logger config)
  (announce-start config)
  (install-signal-traps)
  (cfg/init! config)
  (log/infof "Finalized own configuration")
  (log-if-using-tls config)
  (init-rabbitmq-connection! config)
  (start-http-server config))

(defn shutdown
  []
  (log/infof "Asked to shut down...")
  (System/exit 0))
