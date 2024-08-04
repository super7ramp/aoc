fn hash(input: &str) -> u32 {
    input.chars()
        .map(|c| c as u32)
        .fold(0, |acc, c| ((acc + c) * 17) % 256)
}

fn main() {
    let input = include_str!("../../input.txt");
    let hash_sum = input.split(',').map(hash).sum::<u32>();
    println!("Hash sum: {hash_sum}");
}
