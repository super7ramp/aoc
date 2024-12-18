import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

record Range(long start, long end) {
    boolean contains(final long index) {
        return start <= index && index < end;
    }
}

static final Pattern MUL_OP = Pattern.compile("mul\\((\\d+),(\\d+)\\)");
static final Pattern DO_OP = Pattern.compile("do\\(\\)");
static final Pattern DONT_OP = Pattern.compile("don't\\(\\)");

void main() throws IOException {
    final String input = Files.readString(Path.of("input.txt"));
    part1(input);
    part2(input);
}

void part1(final String input) {
    final Matcher mulMatcher = MUL_OP.matcher(input);
    long total = 0;
    while (mulMatcher.find()) {
        final long operand1 = Long.parseLong(mulMatcher.group(1));
        final long operand2 = Long.parseLong(mulMatcher.group(2));
        total += operand1 * operand2;
    }
    System.out.println("Total part1: " + total);
}

void part2(final String input) {
    final Matcher doMatcher = DO_OP.matcher(input);
    final Matcher dontMatcher = DONT_OP.matcher(input);
    final var dontRanges = new ArrayList<Range>();
    while (dontMatcher.find()) {
        final int dontStart = dontMatcher.start();
        doMatcher.region(dontStart, input.length());
        final int dontEnd = doMatcher.find() ? doMatcher.end() : input.length();
        dontRanges.add(new Range(dontStart, dontEnd));
    }
    System.out.println("Don't ranges: " + dontRanges);

    final Matcher mulMatcher = MUL_OP.matcher(input);
    long total = 0;
    while (mulMatcher.find()) {
        final int mulStart = mulMatcher.start();
        if (dontRanges.stream().noneMatch(range -> range.contains(mulStart))) {
            final long operand1 = Long.parseLong(mulMatcher.group(1));
            final long operand2 = Long.parseLong(mulMatcher.group(2));
            total += operand1 * operand2;
        }
    }
    System.out.println("Total part2: " + total);
}