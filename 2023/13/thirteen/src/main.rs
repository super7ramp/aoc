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

    pub fn score(&self) -> usize {
        let rows_of_reflection = self.find_rows_of_reflections();
        let columns_of_reflection = self.find_columns_of_reflections();

        let row_score = rows_of_reflection.iter()
            .map(|row_above_count| row_above_count * 100)
            .sum::<usize>();

        let column_score = columns_of_reflection.iter().sum::<usize>();

        row_score + column_score
    }

    fn find_rows_of_reflections(&self) -> Vec<usize> {
        (1..self.row_count).filter(|&row_index| {
            let max_height = min(row_index, self.row_count - row_index) + 1;
            (1..max_height).all(|height| {
                let top = self.line_at(row_index - height);
                let bottom = self.line_at(row_index + height - 1);
                top == bottom
            })
        }).collect::<Vec<usize>>()
    }

    fn find_columns_of_reflections(&self) -> Vec<usize> {
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
}

fn main() {
    let patterns = include_str!("../../input.txt").split("\n\n")
        .map(Pattern::from)
        .collect::<Vec<Pattern>>();

    let total_score = patterns.iter()
        .map(Pattern::score)
        .sum::<usize>();
    println!("Total score: {total_score}");
}
