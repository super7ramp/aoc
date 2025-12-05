(ns four
  (:require [clojure.string :as str]))

(defn- parse-grid
  "Parses the input string into a vector of vectors representing the rows and columns."
  [input]
  (->> (str/split-lines input)
       (mapv vec)))

(defn- top
  "The content of the cell directly above the given position."
  [grid [row column]]
  (get-in grid [(dec row) column]))

(defn- top-right
  "The content of the cell diagonally above and to the right of the given position."
  [grid [row column]]
  (get-in grid [(dec row) (inc column)]))

(defn- right
  "The content of the cell directly to the right of the given position."
  [grid [row column]]
  (get-in grid [row (inc column)]))

(defn- bottom-right
  "The content of the cell diagonally below and to the right of the given position."
  [grid [row column]]
  (get-in grid [(inc row) (inc column)]))

(defn- bottom
  "The content of the cell directly below the given position."
  [grid [row column]]
  (get-in grid [(inc row) column]))

(defn- bottom-left
  "The content of the cell diagonally below and to the left of the given position."
  [grid [row column]]
  (get-in grid [(inc row) (dec column)]))

(defn- left
  "The content of the cell directly to the left of the given position."
  [grid [row column]]
  (get-in grid [row (dec column)]))

(defn- top-left
  "The content of the cell diagonally above and to the left of the given position."
  [grid [row column]]
  (get-in grid [(dec row) (dec column)]))

(def
  ^{:private  true
    :arglists '([grid pos])}
  adjacents
  "Returns a vector of the contents of all adjacent cells to the given position."
  (juxt top top-right right bottom-right bottom bottom-left left top-left))

(defn- roll-of-paper?
  "If given a grid and a position, checks if the character at the given position represents a roll of paper (@).
   If given a character, checks if it represents a roll of paper."
  ([grid pos]
   (roll-of-paper? (get-in grid pos)))
  ([character]
   (= \@ character)))

(defn- adjacent-rolls-of-paper-count
  "Counts the number of adjacent rolls of paper around the given position."
  [grid pos]
  (->> (adjacents grid pos)
       (filter roll-of-paper?)
       count))

(defn- accessible?
  "Checks if the roll of paper at the given position is accessible (i.e., has fewer than 4 adjacent rolls of paper)."
  [grid pos]
  (< (adjacent-rolls-of-paper-count grid pos) 4))

(defn accessible-rolls-of-paper
  "Finds all accessible rolls of paper in the given grid."
  [grid]
  (let [row-count (count grid), column-count (count (first grid))]
    (for [row (range row-count), column (range column-count)
          :let [pos [row column]]
          :when (and (roll-of-paper? grid pos) (accessible? grid pos))]
      pos)))

(defn- remove-rolls-of-paper
  "Removes all accessible rolls of paper from the given grid, returning the updated grid and the count of removed rolls."
  [grid]
  (->> (accessible-rolls-of-paper grid)
       (reduce (fn [[current-grid count] pos]
                 (vector (assoc-in current-grid pos \.) (inc count)))
               [grid 0])))

(defn removable-rolls-of-paper-count
  "Counts the total number of rolls of paper that can be removed from the given grid
   by repeatedly removing accessible rolls until none remain."
  [grid]
  (loop [initial-grid grid, removed-total 0]
    (let [[updated-grid, removed] (remove-rolls-of-paper initial-grid)]
      (if (zero? removed)
        removed-total
        (recur updated-grid (+ removed-total removed))))))

(def input (slurp "input.txt"))

(let [grid (parse-grid input)]
  (println "Part 1:")
  (println "Accessible rolls of paper:" (count (accessible-rolls-of-paper grid)))
  (println "Part 2:")
  (println "Removable rolls of paper:" (time (removable-rolls-of-paper-count grid))))
