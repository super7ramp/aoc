import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.*;
import java.util.function.Predicate;

interface Rule extends Predicate<List<Integer>> {}

static class Rules {
    private final Collection<Rule> rules;
    Rules(final Collection<Rule> rules) {
        this.rules = rules;
    }
    static Rules valueOf(final String input) {
        final String[] rulesToParse = input.split("\n");
        final var rules = new ArrayList<Rule>(rulesToParse.length);
        for (final String rule : rulesToParse) {
            final String[] parts = rule.split("\\|");
            final Integer before = Integer.valueOf(parts[0]);
            final Integer after = Integer.valueOf(parts[1]);
            final Rule predicate = pages -> {
                final int afterIndex = pages.indexOf(after);
                return afterIndex == -1 || pages.indexOf(before) < afterIndex;
            };
            rules.add(predicate);
        }
        return new Rules(rules);
    }

    Collection<Rule> all() {
        return Collections.unmodifiableCollection(rules);
    }
}

record Update(List<Integer> pages) {
    static Update valueOf(final String input) {
        final String[] parts = input.split(",");
        final List<Integer> pages = Arrays.stream(parts).map(Integer::valueOf).toList();
        return new Update(pages);
    }
    Integer middlePage() {
        return pages.get(pages.size() / 2);
    }
}

static class SafetyManual {
    private final Rules rules;
    private final List<Update> updates;

    SafetyManual(final Rules rules, final List<Update> updates) {
        this.rules = rules;
        this.updates = updates;
    }

    static SafetyManual valueOf(final String input) {
        final String[] parts = input.split("\n\n");
        final Rules rules = Rules.valueOf(parts[0]);
        final List<Update> updates = Arrays.stream(parts[1].split("\n")).map(Update::valueOf).toList();
        return new SafetyManual(rules, updates);
    }

    List<Update> validatedUpdates() {
        return updates.stream()
                .filter(update -> rules.all().stream().allMatch(rule -> rule.test(update.pages())))
                .toList();
    }
}

void main() throws IOException {
    final String input = Files.readString(Path.of("input.txt"));
    final SafetyManual manual = SafetyManual.valueOf(input);

    final List<Update> validatedUpdates = manual.validatedUpdates();
    System.out.println("Validated updates: " + validatedUpdates);

    final int validatedUpdatesMiddlePageSum = validatedUpdates.stream()
            .mapToInt(Update::middlePage)
            .reduce(0, Integer::sum);
    System.out.println("Sum of validated updates middle pages: " + validatedUpdatesMiddlePageSum);
}