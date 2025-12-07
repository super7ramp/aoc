(ns seven
  (:require [clojure.string :as str]))

(defn- splitter? [content]
  (= \^ content))

(defn- void? [content]
  (= \. content))

(def ^:private beam \|)

(defn- bottom [[row col]]
  [(inc row) col])

(defn- bottom-left [[row col]]
  [(inc row) (dec col)])

(defn- bottom-right [[row col]]
  [(inc row) (inc col)])

(defn parse-manifold
  "Parses the given input into a structure that looks like:

  ```
  {:grid [[. S .]
          [. ^ .]
          [. . .]]
   :starts [[0 1]]
   :splits 0}
  ```
  "
  [input]
  (let [grid (->> (str/split-lines input) (mapv #(mapv identity %)))
        start-column (quot (count (first grid)) 2)]
    {:grid grid, :starts [[0 start-column]], :splits 0}))

(defn propagate-beam
  "Propagates the beam inside the given manifold. Returns the resulting manifold."
  ([manifold]
   (loop [previous-manifold manifold]
     (let [starts (:starts previous-manifold)]
       (if (empty? starts)
         previous-manifold
         (recur (reduce propagate-beam
                        (dissoc previous-manifold :starts)
                        starts))))))

  ([{:keys [grid starts splits] :as manifold} start]
   (let [bottom-pos (bottom start)
         bottom-content (get-in grid bottom-pos)]
     (cond

       (void? bottom-content)
       {:grid   (assoc-in grid bottom-pos beam)
        :starts (conj starts bottom-pos)
        :splits splits}

       (splitter? bottom-content)
       (let [bottom-left-pos (bottom-left start)
             bottom-right-pos (bottom-right start)
             bottom-left-exists (some? (get-in grid bottom-left-pos))
             bottom-right-exists (some? (get-in grid bottom-right-pos))]
         {:grid   (cond-> grid
                    bottom-left-exists (assoc-in bottom-left-pos beam)
                    bottom-right-exists (assoc-in bottom-right-pos beam))
          :starts (cond-> starts
                    bottom-left-exists (conj bottom-left-pos)
                    bottom-right-exists (conj bottom-right-pos))
          :splits (inc splits)})

       :else
       ; outside of grid or existing beam, nothing to add
       manifold))))

(def input (slurp "input.txt"))

(let [manifold (parse-manifold input)]
  (println "Part 1:")
  (println "Total splits:" (:splits (propagate-beam manifold))))