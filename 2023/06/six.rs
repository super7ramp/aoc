fn main() {
    let input = include_str!("input.txt").split("\n").collect::<Vec<&str>>();
    let times = parse_line(&input[0]);
    let distances = parse_line(&input[1]);
    println!("Times: {:?}", times);
    println!("Distances: {:?}", distances);

    let number_of_ways_part_1 = number_of_ways(&times, &distances);
    println!("Number of ways to beat the record (part 1): {number_of_ways_part_1}");

    let number_of_ways_part_2 = number_of_ways(&[58819676], &[434104122191218]);
    println!("Number of ways to beat the record (part 2): {number_of_ways_part_2}");
}

fn parse_line(line: &str) -> Vec<u64>{
    line.split(" ")
        .map(str::trim)
        .filter(|s| !s.is_empty())
        .skip(1)
        .map(|s| s.parse::<u64>().unwrap())
        .collect::<Vec<u64>>()
}

fn number_of_ways(times: &[u64], distances: &[u64]) -> u64 {
    (0..times.len()).into_iter()
        .map(|i| solution_count(times[i], distances[i]))
        .reduce(|a, b| a * b)
        .unwrap()
}

fn solution_count(race_time_in_ms: u64, minimum_distance_in_mm: u64) -> u64 {
    let mut count = 0u64;
    for hold_duration_in_ms in 1..race_time_in_ms {
        let distance_covered_in_ms = (race_time_in_ms - hold_duration_in_ms) * hold_duration_in_ms;
        if distance_covered_in_ms > minimum_distance_in_mm {
            count += 1;
        }
    }
    count
}
