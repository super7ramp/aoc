(ns two
  (:require [clojure.string :as str]))

(defn- parse-range [range-str]
  (let [[start end] (->> (str/split range-str #"-")
                         (map parse-long))]
    (range start (inc end))))

(defn parse-ranges [input]
  (->> (str/split input #",\n?")
       (map parse-range)))

(defn- has-seq-repeated? [string times]
  (if (not (zero? (rem (count string) times)))
    false
    (let [part-length (quot (count string) times)
          parts (->> (range times)
                     (map #(subs string (* % part-length) (* (inc %) part-length))))]
      (apply = parts))))

(defn is-invalid-id? [id]
  (has-seq-repeated? (str id) 2))

(defn is-really-invalid-id? [id]
  (let [id-string (str id), id-string-length (count id-string)]
    (loop [times 2]
      (cond
        (has-seq-repeated? id-string times) true
        (< times id-string-length) (recur (inc times))))))

(def input (slurp "input.txt"))

(let [ranges (flatten (parse-ranges input))]
  (println "Part 1:")
  (println "Invalid ID sum: " (time (->> ranges (filter is-invalid-id?) (reduce +))))
  (println "Part 2:")
  (println "Invalid ID sum: " (time (->> ranges (filter is-really-invalid-id?) (reduce +)))))
