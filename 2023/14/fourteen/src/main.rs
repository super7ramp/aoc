use std::collections::HashMap;
use std::fmt::{Display, Formatter};

#[derive(Eq, Hash, PartialEq, Clone)]
enum Item {
    RoundedRock,
    CubeShapedRock,
    Empty,
}

impl Item {
    fn from_char(c: char) -> Option<Self> {
        match c {
            'O' => Some(Item::RoundedRock),
            '#' => Some(Item::CubeShapedRock),
            '.' => Some(Item::Empty),
            _ => None,
        }
    }
}

impl Display for Item {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        let symbol = match self {
            Item::RoundedRock => 'O',
            Item::CubeShapedRock => '#',
            Item::Empty => '.',
        };
        write!(f, "{}", symbol)
    }
}

struct Platform {
    data: Vec<Item>,
    column_count: usize,
    row_count: usize,
}

impl Platform {
    fn new(data: Vec<Item>, column_count: usize, row_count: usize) -> Self {
        Platform {
            data,
            column_count,
            row_count,
        }
    }

    pub fn from(input: &str) -> Self {
        let rows = input.split("\n").collect::<Vec<&str>>();
        let row_count = rows.len();
        let column_count = rows[0].len();
        let data = rows
            .into_iter()
            .flat_map(|row| row.chars())
            .map(|c| Item::from_char(c).unwrap())
            .collect::<Vec<Item>>();
        Platform::new(data, column_count, row_count)
    }

    pub fn tilt_north(&mut self) {
        for row in 1..self.row_count {
            for column in 0..self.column_count {
                if let Some(Item::RoundedRock) = self.item_at(column, row) {
                    let mut destination_row = row;
                    while let Some(Item::Empty) = self.item_at(column, destination_row - 1) {
                        destination_row -= 1;
                        if destination_row == 0 {
                            break;
                        }
                    }
                    self.swap(column, row, column, destination_row);
                }
            }
        }
    }

    pub fn tilt_west(&mut self) {
        for column in 1..self.column_count {
            for row in 0..self.row_count {
                if let Some(Item::RoundedRock) = self.item_at(column, row) {
                    let mut destination_column = column;
                    while let Some(Item::Empty) = self.item_at(destination_column - 1, row) {
                        destination_column -= 1;
                        if destination_column == 0 {
                            break;
                        }
                    }
                    self.swap(column, row, destination_column, row);
                }
            }
        }
    }

    pub fn tilt_south(&mut self) {
        for row in (0..self.row_count).rev() {
            for column in 0..self.column_count {
                if let Some(Item::RoundedRock) = self.item_at(column, row) {
                    let mut destination_row = row;
                    while let Some(Item::Empty) = self.item_at(column, destination_row + 1) {
                        destination_row += 1;
                    }
                    self.swap(column, row, column, destination_row);
                }
            }
        }
    }

    pub fn tilt_east(&mut self) {
        for column in (0..self.column_count).rev() {
            for row in 0..self.row_count {
                if let Some(Item::RoundedRock) = self.item_at(column, row) {
                    let mut destination_column = column;
                    while let Some(Item::Empty) = self.item_at(destination_column + 1, row) {
                        destination_column += 1;
                    }
                    self.swap(column, row, destination_column, row);
                }
            }
        }
    }

    fn item_at(&self, column: usize, row: usize) -> Option<&Item> {
        self.data.get(self.to_index(column, row))
    }

    fn swap(&mut self, column1: usize, row1: usize, column2: usize, row2: usize) {
        let from = self.to_index(column1, row1);
        let to = self.to_index(column2, row2);
        self.data.swap(from, to);
    }

    fn to_index(&self, column: usize, row: usize) -> usize {
        row * self.column_count + column
    }

    pub fn load(&self) -> usize {
        self.data
            .iter()
            .enumerate()
            .filter(|&(_, item)| matches!(item, Item::RoundedRock))
            .map(|(index, _)| self.row_and_row_below_count(index))
            .sum()
    }

    fn row_and_row_below_count(&self, index: usize) -> usize {
        self.row_count - self.row_number(index)
    }

    fn row_number(&self, index: usize) -> usize {
        index / self.column_count
    }

    pub fn n_tilt_cycles(&mut self, n: usize) {
        let mut tilt_cache = HashMap::new();
        let mut i = 0;
        while i < n {
            if let Some(previous_occurrence) = tilt_cache.get(&self.data) {
                let loop_length = i - previous_occurrence;
                println!("Detected same occurrence at cycle {i} as at cycle {previous_occurrence}");
                println!("Loop is {loop_length}-long");
                i = n - ((n - previous_occurrence) % loop_length);
                println!("Short-circuited to cycle {}", i);
                tilt_cache.clear();
                continue;
            }
            tilt_cache.insert(self.data.clone(), i);
            self.tilt_north();
            self.tilt_west();
            self.tilt_south();
            self.tilt_east();
            i += 1;
        }
    }
}

impl Display for Platform {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        for (i, item) in self.data.iter().enumerate() {
            item.fmt(f)?;
            if i % self.column_count == self.column_count - 1 {
                writeln!(f)?;
            }
        }
        Ok(())
    }
}

fn main() {
    let mut platform = Platform::from(include_str!("../../input-example.txt"));
    println!("{platform}");
    println!("Initial load: {}", platform.load());

    println!("Tilting north");
    platform.tilt_north();
    println!("{platform}");
    println!("Load: {}", platform.load());

    println!("Tilting west");
    platform.tilt_west();
    println!("{platform}");
    
    println!("Tilting south");
    platform.tilt_south();
    println!("{platform}");

    println!("Tilting east");
    platform.tilt_east();
    println!("{platform}");

    println!("Load: {}", platform.load());

    println!("Performing 2 other cycles");
    platform.n_tilt_cycles(2);
    println!("{platform}");
    println!("Load: {}", platform.load());

    println!("Performing 999 999 997 other cycles");
    platform.n_tilt_cycles(999_999_997);
    println!("{platform}");
    println!("Load: {}", platform.load());
}
