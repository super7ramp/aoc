use std::cmp::min;

struct Pattern {
    pattern: String,
    column_count: usize,
    row_count: usize,
}

impl Pattern {
    fn new(pattern: String, column_count: usize, row_count: usize) -> Pattern {
        Pattern {
            pattern,
            column_count,
            row_count,
        }
    }

    fn from(input: &str) -> Pattern {
        let rows = input.split("\n").collect::<Vec<&str>>();
        let row_count = rows.len();
        let column_count = rows[0].len();
        let pattern = input.to_string();
        Pattern::new(pattern, column_count, row_count)
    }

    pub fn score(&self) -> (usize, usize) {
        let reflection_rows = self.reflection_rows();
        let reflection_columns = self.reflection_columns();
        let row_score = reflection_rows.iter()
            .map(|row_above_count| row_above_count * 100)
            .sum::<usize>();
        let column_score = reflection_columns.iter().sum::<usize>();
        let score_part_1 = row_score + column_score;

        let (smudged_reflection_rows, smudged_reflection_columns) =
            self.smudged_reflections(reflection_rows, reflection_columns);
        let smudged_row_score = smudged_reflection_rows.iter()
            .map(|row_above_count| row_above_count * 100)
            .sum::<usize>();
        let smudged_column_score = smudged_reflection_columns.iter().sum::<usize>();
        let score_part_2 = smudged_row_score + smudged_column_score;

        (score_part_1, score_part_2)
    }

    fn reflection_rows(&self) -> Vec<usize> {
        (1..self.row_count).filter(|&row_index| {
            let max_height = min(row_index, self.row_count - row_index) + 1;
            (1..max_height).all(|height| {
                let top = self.line_at(row_index - height);
                let bottom = self.line_at(row_index + height - 1);
                top == bottom
            })
        }).collect::<Vec<usize>>()
    }

    fn reflection_columns(&self) -> Vec<usize> {
        (1..self.column_count).filter(|&column_index| {
            let width = min(column_index, self.column_count - column_index);
            (0..self.row_count).all(|row_index| {
                let left = &self.line_at(row_index)[column_index - width..column_index];
                let reverted_right = &self.line_at(row_index)[column_index..column_index + width]
                    .chars()
                    .rev()
                    .collect::<String>();
                left == reverted_right
            })
        }).collect::<Vec<usize>>()
    }

    fn line_at(&self, row: usize) -> &str {
        &self.pattern[self.index_of(row, 0)..self.index_of(row, self.column_count)]
    }

    fn index_of(&self, row: usize, column: usize) -> usize {
        row * (self.column_count + 1 /* newline */) + column
    }

    fn smudged_reflections(&self, original_rows: Vec<usize>, original_columns: Vec<usize>) -> (Vec<usize>, Vec<usize>) {
        for row_index in 0..self.row_count {
            for column_index in 0..self.column_count {
                let mut new_pattern = self.pattern.clone();
                let index = self.index_of(row_index, column_index);
                let character = new_pattern.chars().nth(index).unwrap();
                if character == '#' {
                    new_pattern.replace_range(index..index + 1, ".");
                } else {
                    new_pattern.replace_range(index..index + 1, "#");
                }

                let smudged_pattern = Pattern::new(new_pattern, self.column_count, self.row_count);
                let rows = smudged_pattern.reflection_rows()
                    .into_iter()
                    .filter(|row| !original_rows.contains(row))
                    .collect::<Vec<usize>>();
                let columns = smudged_pattern.reflection_columns()
                    .into_iter()
                    .filter(|row| !original_columns.contains(row))
                    .collect::<Vec<usize>>();

                if !rows.is_empty() || !columns.is_empty() {
                    return (rows, columns);
                }
            }
        }
        return (vec![], vec![]);
    }
}

fn main() {
    let patterns = include_str!("input.txt").split("\n\n")
        .map(Pattern::from)
        .collect::<Vec<Pattern>>();

    let (total_score_part_1, total_score_part_2) = patterns.iter()
        .map(Pattern::score)
        .reduce(|(acc_row, acc_column), (row, column)| (acc_row + row, acc_column + column))
        .unwrap();
    println!("Total score (part 1): {total_score_part_1}");
    println!("Total score (part 2): {total_score_part_2}");
}
