(ns five
  (:require [clojure.string :as str]))

(defn- parse-id-range
  "Parses a string of the form 'start-end' into a map that looks like `{:start 3, :end 5}`."
  [string]
  (let [[start end] (->> (str/split string #"-") (map parse-long))]
    {:start start :end end}))

(defn- parse-db
  "Parses the input into a db structure that looks like:

  ```
  {:fresh-id-ranges [{:start 3, :end 5}
                     {:start 10, :end 14}
                     {:start 16, :end 20}
                     {:start 12, :end 18}]
   :ingredient-ids [1 5 8 11 17 32]}
  ```
  "
  [input]
  (let [[first-section second-section] (str/split input #"\n\n")
        fresh-id-ranges (->> (str/split-lines first-section) (mapv parse-id-range))
        ingredient-ids (->> (str/split-lines second-section) (mapv parse-long))]
    {:fresh-id-ranges fresh-id-ranges
     :ingredient-ids  ingredient-ids}))

(defn- includes?
  "Returns true iff the given range includes the given id."
  [range id]
  (and (some? range)
       (>= id (:start range))
       (<= id (:end range))))

(defn- is-fresh?
  "Returns true iff the given id is in any of the fresh id ranges."
  [db id]
  (->> (:fresh-id-ranges db)
       (some #(includes? % id))))

(defn existing-fresh-ingredient-ids
  "Returns the existing ingredient ids that are fresh."
  [db]
  (->> (:ingredient-ids db)
       (filter #(is-fresh? db %))))

(defn- range-reducer
  "A reducing function that merges ranges that overlap. It assumes that given ranges are sorted by :start."
  [ranges new-range]
  (let [last-range (peek ranges)]
    (cond
      (not (includes? last-range (:start new-range))) (conj ranges new-range)
      (includes? last-range (:end new-range)) ranges
      :else (let [last-index (dec (count ranges))
                  updated-last-range (assoc last-range :end (:end new-range))]
              (assoc ranges last-index updated-last-range)))))

(defn- range-size
  "Returns the size of the given range."
  [range]
  (inc (- (:end range) (:start range))))

(defn count-fresh-ingredient-ids
  "Returns the total count of fresh ingredient ids."
  [db]
  (->> (:fresh-id-ranges db)
       (sort-by :start)
       (reduce range-reducer [])
       (map range-size)
       (reduce +)))

(def input (slurp "input.txt"))

(let [db (parse-db input)]
  (println "Part 1:")
  (println "Existing fresh ingredients:" (count (existing-fresh-ingredient-ids db)))
  (println "Part 2:")
  (println "Fresh ingredients:" (count-fresh-ingredient-ids db)))