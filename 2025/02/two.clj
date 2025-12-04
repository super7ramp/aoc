(ns two
  (:require [clojure.string :as str]))

(defn- parse-range [range-string]
  (let [[start end] (->> (str/split range-string #"-")
                         (map parse-long))]
    (range start (inc end))))

(defn parse-ranges [input]
  (->> (str/split input #",\n?")
       (map parse-range)))

(defn- has-repeated-substring? [string times]
  (let [string-length (count string)]
    (if (not (zero? (rem string-length times)))
      false
      (let [substring-length (quot string-length times)
            substrings (->> (range times)
                            (map #(subs string
                                        (* % substring-length)
                                        (* (inc %) substring-length))))]
        (apply = substrings)))))

(defn is-invalid-id? [id]
  (has-repeated-substring? (str id) 2))

(defn is-really-invalid-id? [id]
  (let [id-string (str id), id-string-length (count id-string)]
    (loop [times 2]
      (cond
        (has-repeated-substring? id-string times) true
        (< times id-string-length) (recur (inc times))
        :else false))))

(def input (slurp "input.txt"))

(let [ids (flatten (parse-ranges input))]
  (println "Part 1:")
  (println "Invalid ID sum: " (time (->> ids (filter is-invalid-id?) (reduce +))))
  (println "Part 2:")
  (println "Invalid ID sum: " (time (->> ids (filter is-really-invalid-id?) (reduce +)))))
