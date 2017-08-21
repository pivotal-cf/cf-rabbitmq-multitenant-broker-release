(ns io.pivotal.pcf.rabbitmq.server-test
  (:require [clojure.test :refer :all]
            [taoensso.timbre :as log]
            [io.pivotal.pcf.rabbitmq.resources :as rs]
            [io.pivotal.pcf.rabbitmq.server :as server]))

(defn LogHolder
  [message & args]
  (def LastLogMessage (list message args)))

(defn NoOp
  [& args]
  (print "NoOp got called with args: " args))

(deftest create-service
  (testing "should log on service creation"
    (with-redefs [server/log-info LogHolder rs/grant-broker-administrator-permissions NoOp]
      (def LastLogMessage nil)
      (server/create-service {:params {:id "my service"}})
      (is (= LastLogMessage  (list "Asked to provision a service: %s" (list "my service"))))
      )))
