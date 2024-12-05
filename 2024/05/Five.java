import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.*;
import java.util.function.Predicate;

record Rule(Integer before, Integer after) implements Predicate<List<Integer>> {
    @Override
    public boolean test(final List<Integer> pages) {
        final int afterIndex = pages.indexOf(after);
        return afterIndex == -1 || pages.indexOf(before) < afterIndex;
    }
}

static class Rules {
    private final Collection<Rule> rules;

    Rules(final Collection<Rule> rules) {
        this.rules = rules;
    }

    static Rules valueOf(final String input) {
        final String[] rulesToParse = input.split("\n");
        final var rules = new ArrayList<Rule>(rulesToParse.length);
        for (final String ruleToParse : rulesToParse) {
            final String[] parts = ruleToParse.split("\\|");
            final Integer before = Integer.valueOf(parts[0]);
            final Integer after = Integer.valueOf(parts[1]);
            final Rule rule = new Rule(before, after);
            rules.add(rule);
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

record UpdateWithInfringedRules(Update update, List<Rule> infringedRules) {
    Update correctedUpdate() {
        if (infringedRules.isEmpty()) {
            return update;
        }
        final List<Integer> correctedPages = new ArrayList<>(update.pages());
        for (final Rule rule : infringedRules) {
            final int beforeIndex = correctedPages.indexOf(rule.before);
            final int afterIndex = correctedPages.indexOf(rule.after);
            if (beforeIndex > afterIndex) {
                correctedPages.remove(beforeIndex);
                correctedPages.add(afterIndex, rule.before());
            }
        }
        return new Update(correctedPages);
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

    List<Update> correctedUpdates() {
        final List<Update> nonValidUpdates = updates.stream()
                .filter(update -> rules.all().stream().anyMatch(rule -> !rule.test(update.pages())))
                .toList();
        return correct(nonValidUpdates);
    }

    private List<Update> correct(final List<Update> updatesToCorrect) {
        final List<UpdateWithInfringedRules> updatesWithInfringedRules = updatesToCorrect.stream()
                .map(update -> {
                    final List<Rule> infringedRules = rules.all().stream()
                            .filter(rule -> !rule.test(update.pages()))
                            .toList();
                    return new UpdateWithInfringedRules(update, infringedRules);
                })
                .toList();

        final List<Update> correctedUpdates = updatesWithInfringedRules.stream()
                .map(UpdateWithInfringedRules::correctedUpdate)
                .toList();

        if (correctedUpdates.stream().anyMatch(correctedUpdate -> rules.all().stream().anyMatch(rule -> !rule.test(correctedUpdate.pages())))) {
            System.out.println("Correction is incomplete, iterating");
            return correct(correctedUpdates);
        }

        return correctedUpdates;
    }
}

void main() throws IOException {
    final String input = Files.readString(Path.of("input.txt"));
    final SafetyManual manual = SafetyManual.valueOf(input);

    final List<Update> validatedUpdates = manual.validatedUpdates();
    System.out.println("(Part 1) Validated updates: " + validatedUpdates);

    final int validatedUpdatesMiddlePageSum = validatedUpdates.stream().mapToInt(Update::middlePage).sum();
    System.out.println("(Part 1) Sum of validated updates middle pages: " + validatedUpdatesMiddlePageSum);

    final List<Update> correctedUpdates = manual.correctedUpdates();
    System.out.println("(Part 2) Corrected updates: " + correctedUpdates);
    final int correctedUpdatesMiddlePageSum = correctedUpdates.stream().mapToInt(Update::middlePage).sum();
    System.out.println("(Part 2) Sum of corrected updates middle pages: " + correctedUpdatesMiddlePageSum);
}