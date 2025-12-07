(ns six
  (:require [clojure.string :as str]))

(defn- find-blank-regions
  "
  Returns the regions of the given string that contains blanks, structured like this:

  ```
  ; input: \"123 328  51 64\"
  [{:start 3, :end 4}
   {:start 7, :end 9}
   {:start 11, :end 12}]
  ```
  "
  [string]
  (let [matcher (re-matcher #"\s+" string)]
    (loop [regions []]
      (if-not (re-find matcher)
        regions
        (let [region {:start (.start matcher) :end (.end matcher)}]
          (recur (conj regions region)))))))

(defn- remove-zero-starts
  "Removes first region if it starts a zero."
  [blank-regions]
  (if (zero? (:start (first blank-regions)))
    (rest blank-regions)
    blank-regions))

(defn- intersect
  "Returns the intersection of the two given regions."
  [x y]
  (if-not (and (some? x) (some? y))
    nil
    (let [{a :start b :end} x
          {c :start d :end} y
          start (max a c)
          end (min b d)]
      (when (< start end)
        {:start start :end end}))))

(defn- intersecting
  "Returns a reducing function that merges intersecting regions."
  [disjoint-regions region]
  (let [previous-region (peek disjoint-regions)
        intersection (intersect previous-region region)]
    (if-not intersection
      (conj disjoint-regions region)
      (assoc disjoint-regions (dec (count disjoint-regions)) intersection))))

(defn- by-intersecting
  "Returns a reducing function that merge a sequence of region vectors into a single region vector containing the
   intersections between all region vectors."
  [previous current]
  (->> (interleave previous current)
       (sort-by :start)
       (reduce intersecting [])))

(defn- find-separator-regions
  "Finds blank separator regions in the given lines."
  [lines]
  (->> lines
       (map find-blank-regions)
       (map remove-zero-starts)
       (reduce by-intersecting)))

(defn- columns
  "Returns the columns separated by blanks in the given input.

  Note that it's not a simple trim: Blanks within a colum are preserved, e.g. called with:

  ```
  123 328  51 64
   45 64  387 23
    6 98  215 314
  *   +   *   +
  ```

  The function returns:

  ```
    [[\"123\" \" 45\" \"  6\" \"*  \"]
     [\"328\" \"64 \" \"98 \" \"+  \"]
     [\" 51\" \"387\" \"215\" \"*  \"]
     [\"64\" \"23\" \"314\" \"+\"]]
  ```
  "
  [string]
  (let [lines (str/split-lines string)]
    (loop [columns []
           separators (find-separator-regions lines)
           start 0]
      (let [next-separator (first separators)
            end (:start next-separator)
            column (->> lines (mapv #(subs % start (or end (count %)))))
            updated-columns (conj columns column)]
        (if (nil? next-separator)
          updated-columns
          (recur updated-columns
                 (rest separators)
                 (:end next-separator)))))))

(defn- parse-operator
  "Parses an operator string into the corresponding function."
  [string]
  (let [trimmed (str/trim string)]
    (cond
      (= "+" trimmed) +
      (= "*" trimmed) *)))

(defn- parse-operands
  "Parses operand rows into actual operands."
  [rows {writing-mode :writing-mode}]
  (if (= :vertical writing-mode)
    (let [width (apply max (map count rows))]
      (->> (range width)
           (mapv (fn [char-index]
                   (->> rows
                        (filter #(< char-index (count %)))
                        (map #(subs % char-index (inc char-index)))
                        (remove str/blank?)
                        (apply str)
                        parse-long)))))
    (->> rows
         (map str/trim)
         (map parse-long))))

(defn- parse-problem
  "Parses a single problem."
  [rows opts]
  {:operator (parse-operator (last rows))
   :operands (parse-operands (butlast rows) opts)})

(defn parse-problems
  "Parses the problem into a map that looks like:

  ```
  [{:operator #object [clojure.core$_STAR_...]
    :operands [1 24 356]}
   {:operator #object [clojure.core$_PLUS_...]
    :operands [369 248 8]}
   {:operator #object [clojure.core$_STAR_...]
    :operands [32 581 175]}
   {:operator #object [clojure.core$_PLUS_...]
    :operands [623 431 4]}]
  ```

  Supports two `:writing-mode`s:

  - `:horizontal`: Operands are parsed horizontally, left-to-right.
  - `:vertical`: Operands are parsed vertically, top-to-bottom.
  "
  ([input]
   (parse-problems input {:writing-mode :horizontal}))
  ([input opts]
   (->> (columns input)
        (mapv #(parse-problem % opts)))))

(defn- solve
  "Solves a problem."
  [{operator :operator, operands :operands}]
  (apply operator operands))

(defn grand-total
  "Gives the grand total, i.e. the sum of all problem results."
  [problems]
  (->> problems
       (map solve)
       (reduce +)))

(def input (slurp "input.txt"))

(let [problems (parse-problems input)]
  (println "Part 1:")
  (println "Grand total:" (grand-total problems)))

(let [problems (parse-problems input {:writing-mode :vertical})]
  (println "Part 2:")
  (println "Grand total:" (grand-total problems)))
