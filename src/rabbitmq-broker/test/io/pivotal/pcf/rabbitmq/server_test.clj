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

(deftest create-service
  (testing "should log on service creation"
    (with-redefs [rs/grant-broker-administrator-permissions NoOp]
      (io/delete-file logfile true)
      (server/create-service {:params {:id "my service"}})
      (is (.contains (get-logs) "Asked to provision a service: my service"))))
  (testing "should log error when rabbitmq is down during service creation"
    (with-redefs [rs/vhost-exists? ThrowException]
      (io/delete-file logfile true)
      (server/create-service {:params {:id "my service"}})
      (is (.contains (get-logs) "Failed to provision a service: my service")))))

(deftest delete-service
  (testing "should log on service deletion"
    (with-redefs [rs/vhost-exists? NoOp]
      (io/delete-file logfile true)
      (server/delete-service {:params {:id "my service"}})
      (is (.contains (get-logs) "Asked to deprovision a service: my service"))))
  (testing "should log error when rabbitmq is down during service deletion"
    (with-redefs [rs/vhost-exists? ThrowException]
      (io/delete-file logfile true)
      (server/delete-service {:params {:id "my service"}})
      (is (.contains (get-logs) "Failed to deprovision a service: my service")))))

(deftest bind-service
  (testing "should log on service binding"
    (io/delete-file logfile true)
    (server/bind-service {:params {:instance_id "my service"}})
    (is (.contains (get-logs) "Asked to bind a service: my service")))
  (testing "should log error when rabbitmq is down on service binding"
    (with-redefs [rs/vhost-exists? ThrowException]
      (io/delete-file logfile true)
      (server/bind-service {:params {:instance_id "my service" :id "user id"}})
      (is (.contains (get-logs) "Failed to bind a service: my service")))))
