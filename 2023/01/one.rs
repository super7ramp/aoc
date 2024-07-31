fn main() {
    let input = include_str!("input.txt");
    let sum = input.split("\n")
        .map(prepare_line)
        .map(find_calibration_value)
        .sum::<u32>();
    println!("{sum}");
}

fn prepare_line(line: &str) -> String {
    line.replace("one", "o1e")
        .replace("two", "t2o")
        .replace("three", "t3e")
        .replace("four", "4")
        .replace("five", "5e")
        .replace("six", "6")
        .replace("seven", "7n")
        .replace("eight", "8t")
        .replace("nine", "n9e")
}

fn find_calibration_value(line: String) -> u32 {
    let digits = line.chars()
        .filter(|c| c.is_ascii_digit())
        .map(|c| c.to_digit(10).unwrap())
        .collect::<Vec<u32>>();
    let calibration_value = digits[0] * 10 + digits.last().unwrap();
    println!("{line} -> {calibration_value}");
    calibration_value
}