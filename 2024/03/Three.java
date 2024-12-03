import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

static final Pattern MUL_OP = Pattern.compile("mul\\((\\d+),(\\d+)\\)");

void main() throws IOException {
    final String input = Files.readString(Path.of("input.txt"));
    final Matcher matcher = MUL_OP.matcher(input);
    long total = 0;
    while (matcher.find()) {
        final long operand1 = Long.parseLong(matcher.group(1));
        final long operand2 = Long.parseLong(matcher.group(2));
        total += operand1 * operand2;
    }
    System.out.println("Total: " + total);
}
