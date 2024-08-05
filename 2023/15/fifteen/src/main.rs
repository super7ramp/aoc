use std::array;
use std::fmt::{Display, Formatter};

#[derive(Clone)]
struct Lens {
    label: String,
    focal_length: u32,
}

trait Instruction {
    fn eval(&self, boxes: &mut [Vec<Lens>; 256]);
}

struct Put {
    lens: Lens,
}

impl Instruction for Put {
    fn eval(&self, boxes: &mut [Vec<Lens>; 256]) {
        let box_number = hash(&self.lens.label) as usize;
        let lenses = boxes.get_mut(box_number).unwrap();
        if let Some(insert_position) = lenses.iter().position(|lens| lens.label == self.lens.label)
        {
            lenses[insert_position] = self.lens.clone();
        } else {
            lenses.push(self.lens.clone());
        }
    }
}

struct Delete {
    label: String,
}

impl Instruction for Delete {
    fn eval(&self, boxes: &mut [Vec<Lens>; 256]) {
        let box_number = hash(&self.label) as usize;
        let lenses = boxes.get_mut(box_number).unwrap();
        if let Some(position) = lenses.iter().position(|lens| lens.label == self.label) {
            lenses.remove(position);
        }
    }
}

struct InstructionSequence {
    instructions: Vec<Box<dyn Instruction>>,
}

impl InstructionSequence {
    fn parse(input: &str) -> Self {
        let instructions = input
            .split(',')
            .map(|instruction| {
                if instruction.ends_with('-') {
                    Self::parse_delete(instruction)
                } else {
                    Self::parse_put(instruction)
                }
            })
            .collect();
        Self { instructions }
    }

    fn parse_delete(input: &str) -> Box<dyn Instruction> {
        let label = input.chars().take(input.len() - 1).collect::<String>();
        Box::new(Delete { label })
    }

    fn parse_put(input: &str) -> Box<dyn Instruction> {
        let instruction_parts = input.split('=').collect::<Vec<&str>>();
        let label = instruction_parts[0].to_string();
        let focal_length = instruction_parts[1].parse().unwrap();
        let lens = Lens {
            label,
            focal_length,
        };
        Box::new(Put { lens })
    }

    fn to_slice(&self) -> &[Box<dyn Instruction>] {
        &self.instructions
    }
}

struct Boxes {
    boxes: [Vec<Lens>; 256],
}

impl Boxes {
    fn new() -> Self {
        let boxes: [Vec<Lens>; 256] = array::from_fn(|_| Vec::new());
        Self { boxes }
    }

    pub fn process(&mut self, initialization_sequence: &[Box<dyn Instruction>]) {
        for instruction in initialization_sequence.iter() {
            instruction.eval(&mut self.boxes);
        }
    }

    pub fn focusing_power(&self) -> u32 {
        self.boxes
            .iter()
            .enumerate()
            .map(|(box_number, lenses)| {
                lenses
                    .iter()
                    .enumerate()
                    .map(|(slot_number, lens)| {
                        Self::lens_focusing_power(box_number, slot_number, lens.focal_length)
                    })
                    .sum::<u32>()
            })
            .sum()
    }

    fn lens_focusing_power(box_number: usize, slot_number: usize, focal_length: u32) -> u32 {
        (1 + box_number as u32) * (1 + slot_number as u32) * focal_length
    }
}

impl Display for Boxes {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        for (box_number, lenses) in self.boxes.iter().enumerate() {
            if lenses.is_empty() {
                continue;
            }
            write!(f, "Box {}:", box_number)?;
            for lens in lenses {
                write!(f, " [{} {}]", lens.label, lens.focal_length)?;
            }
            writeln!(f)?;
        }
        Ok(())
    }
}

fn hash(input: &str) -> u32 {
    input
        .chars()
        .map(|c| c as u32)
        .fold(0, |acc, c| ((acc + c) * 17) % 256)
}

fn main() {
    let input = include_str!("../../input.txt");
    let hash_sum = input.split(',').map(hash).sum::<u32>();
    println!("Hash sum: {hash_sum}");

    let instructions = InstructionSequence::parse(input);
    let mut boxes = Boxes::new();
    boxes.process(instructions.to_slice());
    println!("{boxes}");

    let focusing_power = boxes.focusing_power();
    println!("Focusing power: {focusing_power}")
}
