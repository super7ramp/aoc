(ns one
  (:require [clojure.string :as str]))

(defn- parse-rotation [line]
  (* (if (= \L (first line)) -1 1)
     (parse-long (subs line 1))))

(defn parse-rotations [input]
  (->> (str/split-lines input)
       (map parse-rotation)))

(defn count-rotations-pointed-at-zero [rotations initial-orientation]
  (->> rotations
       (reduce (fn [[orientation count] rotation]
                 (let [new-orientation (mod (+ orientation rotation) 100),
                       new-count (if (== new-orientation 0) (inc count) count)]
                   [new-orientation new-count]))
               [initial-orientation 0])
       second))

(defn- count-rotation-crossed-zero [rotation initial-orientation]

  (let [crossed-zero-count (quot (abs rotation) 100)
        rotation-leftover (rem rotation 100)
        unbound-orientation (+ initial-orientation rotation-leftover)]

    (if (or (and (<= unbound-orientation 0) (not (zero? initial-orientation)))
            (>= unbound-orientation 100))
      (inc crossed-zero-count)
      crossed-zero-count)))

(defn count-rotations-crossed-zero [rotations initial-orientation]
  (->> rotations
       (reduce (fn [[orientation count] rotation]
                 (let [new-orientation (mod (+ orientation rotation) 100),
                       new-count (+ count (count-rotation-crossed-zero rotation orientation))]
                   [new-orientation new-count]))
               [initial-orientation 0])
       second))

(def input (slurp "input.txt"))

(let [rotations (parse-rotations input), initial-rotation 50]
  (println "Part 1:")
  (printf "Pointed at orientation zero %d times\n" (count-rotations-pointed-at-zero rotations initial-rotation))
  (println "Part 2:")
  (printf "Crossed orientation zero %d times\n" (count-rotations-crossed-zero rotations initial-rotation)))
