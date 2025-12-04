(ns three
  (:require [clojure.string :as str]))

(defn- parse-bank [line]
  (->> line
       (mapv #(parse-long (str %)))))

(defn parse-banks [input]
  (->> (str/split-lines input)
       (map parse-bank)))

(defn- comparing [{by :by, reversed? :reversed?, then-compare :then}]
  (fn [a b]
    (let [comparison (if reversed?
                       (compare (by b) (by a))
                       (compare (by a) (by b)))]
      (if (and (zero? comparison) (some? then-compare))
        (then-compare a b)
        comparison))))

(defrecord IndexedBattery [index joltage])

(def ^:private comparing-battery-by-joltage-then-by-index-reversed
  (comparing {:by :joltage :then (comparing {:by :index :reversed? true})}))

(defn- pow [base exp]
  (reduce * (repeat exp base)))

(defn- find-max-joltage [batteries turned-on-battery-count]
  (let [battery-count (count batteries)]
    (loop [selected-count 0, joltage 0, search-start 0]
      (if (== selected-count turned-on-battery-count)
        joltage
        (let [search-end (- battery-count (- turned-on-battery-count selected-count 1))
              joltest-battery (->> (subvec batteries search-start search-end)
                                   (map-indexed ->IndexedBattery)
                                   (sort comparing-battery-by-joltage-then-by-index-reversed)
                                   last)
              added-joltage (* (:joltage joltest-battery) (pow 10 (- turned-on-battery-count selected-count 1)))]
          (recur (inc selected-count)
                 (+ joltage added-joltage)
                 (+ search-start (:index joltest-battery) 1)))))))

(defn find-max-total-joltage
  ([banks]
   (find-max-total-joltage banks 2))
  ([banks turned-on-battery-count]
   (->> banks
        (map #(find-max-joltage % turned-on-battery-count))
        (reduce +))))

(def input (slurp "input.txt"))

(let [banks (parse-banks input)]
  (println "Part 1:")
  (println "Max total joltage (2 batteries turned on by bank):" (find-max-total-joltage banks))
  (println "Part 2:")
  (println "Max total joltage (12 batteries turned on by bank):" (find-max-total-joltage banks 12)))
