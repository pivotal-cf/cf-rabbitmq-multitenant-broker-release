(ns io.pivotal.pcf.rabbitmq.server-test
  (:require [clojure.test :refer :all]
            [clojure.java.io :as io]
            [clojure.string :as string]
            [taoensso.timbre :as log]
            [io.pivotal.pcf.rabbitmq.resources :as rs]
            [io.pivotal.pcf.rabbitmq.server :as server]))

(def logfile "target/create-service.log")

(defn turn-on-logging
  [run-test]
  (log/set-config! [:appenders :spit :enabled?] true)
  (log/set-config! [:shared-appender-config :spit-filename] logfile)
  (run-test))

(use-fixtures :once turn-on-logging)

(defn NoOp
  [& args]
  ())

(defn ThrowException
  [vhost]
  (throw (Exception. "Throw exception explicitly to test")))

(defn get-logs
  []
  (with-open [stream (io/reader logfile)]
    (format "%s\n" (string/join "\n" (line-seq stream)))))
