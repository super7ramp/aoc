use std::fmt::{Display, Formatter};

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

    fn row_number(&self, index: usize) -> usize {
        index / self.column_count
    }

    fn row_and_row_below_count(&self, index: usize) -> usize {
        self.row_count - self.row_number(index)
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
    let mut platform = Platform::from(include_str!("../../input.txt"));
    println!("{platform}");

    println!("Tilting north");
    platform.tilt_north();
    println!("{platform}");

    println!("Load: {}", platform.load());
}
