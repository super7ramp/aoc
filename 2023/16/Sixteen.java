import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.*;
import java.util.stream.Stream;

import static java.util.function.Predicate.not;

enum Direction {
    NORTH,
    SOUTH,
    EAST,
    WEST;
}

enum Element {
    EMPTY('.', deflectionMapOf(
            Direction.NORTH, EnumSet.of(Direction.NORTH),
            Direction.EAST, EnumSet.of(Direction.EAST),
            Direction.SOUTH, EnumSet.of(Direction.SOUTH),
            Direction.WEST, EnumSet.of(Direction.WEST))),
    SPLITTER_NS('|', deflectionMapOf(
            Direction.NORTH, EnumSet.of(Direction.NORTH),
            Direction.EAST, EnumSet.of(Direction.NORTH, Direction.SOUTH),
            Direction.SOUTH, EnumSet.of(Direction.SOUTH),
            Direction.WEST, EnumSet.of(Direction.NORTH, Direction.SOUTH)
    )),
    SPLITTER_WE('-', deflectionMapOf(
            Direction.NORTH, EnumSet.of(Direction.WEST, Direction.EAST),
            Direction.EAST, EnumSet.of(Direction.EAST),
            Direction.SOUTH, EnumSet.of(Direction.WEST, Direction.EAST),
            Direction.WEST, EnumSet.of(Direction.WEST)
    )),
    MIRROR_ES_NW('\\', deflectionMapOf(
            Direction.NORTH, EnumSet.of(Direction.WEST),
            Direction.EAST, EnumSet.of(Direction.SOUTH),
            Direction.SOUTH, EnumSet.of(Direction.EAST),
            Direction.WEST, EnumSet.of(Direction.NORTH)
    )),
    MIRROR_EN_SW('/', deflectionMapOf(
            Direction.NORTH, EnumSet.of(Direction.EAST),
            Direction.EAST, EnumSet.of(Direction.NORTH),
            Direction.SOUTH, EnumSet.of(Direction.WEST),
            Direction.WEST, EnumSet.of(Direction.SOUTH)
    ));

    private static final Element[] CACHED_VALUES = values();

    private final Map<Direction, Set<Direction>> deflections;

    private final char symbol;

    Element(final char symbol, final Map<Direction, Set<Direction>> deflections) {
        this.symbol = symbol;
        this.deflections = deflections;
    }

    private static <V> Map<Direction, V> deflectionMapOf(
            final Direction k1, final V v1,
            final Direction k2, final V v2,
            final Direction k3, final V v3,
            final Direction k4, final V v4) {
        final var map = new EnumMap<Direction, V>(Direction.class);
        map.put(k1, v1);
        map.put(k2, v2);
        map.put(k3, v3);
        map.put(k4, v4);
        return map;
    }

    static Element valueOf(final char symbol) {
        return Arrays.stream(CACHED_VALUES)
                .filter(element -> element.symbol == symbol)
                .findFirst()
                .orElseThrow();
    }

    Set<Direction> deflectBeamGoingTo(final Direction direction) {
        return deflections.get(direction);
    }
}

record Position(int row, int column) {
    Position to(final Direction direction) {
        return switch (direction) {
            case NORTH -> new Position(row - 1, column);
            case EAST -> new Position(row, column + 1);
            case SOUTH -> new Position(row + 1, column);
            case WEST -> new Position(row, column - 1);
        };
    }
}

static class Contraption {

    class BeamPart {
        private final Position position;
        private final Direction direction;

        BeamPart(final Position position, final Direction direction) {
            this.position = position;
            this.direction = direction;
        }

        Stream<BeamPart> nextParts() {
            return elementAt(position)
                    .deflectBeamGoingTo(direction).stream()
                    .map(newDirection -> new BeamPart(position.to(newDirection), newDirection))
                    .filter(beamPart -> contains(beamPart.position));
        }

        @Override
        public boolean equals(final Object obj) {
            if (this == obj) {
                return true;
            }
            if (!(obj instanceof BeamPart other)) {
                return false;
            }
            return position.equals(other.position) && direction == other.direction;
        }

        @Override
        public int hashCode() {
            return Objects.hash(position, direction);
        }
    }

    private final Element[][] elements;

    private Contraption(final Element[][] elements) {
        this.elements = elements;
    }

    static Contraption valueOf(final String input) {
        final var lines = input.lines().toList();
        final var elements = new Element[lines.size()][];
        for (int i = 0; i < lines.size(); i++) {
            elements[i] = lines.get(i).chars()
                    .mapToObj(c -> Element.valueOf(Character.toChars(c)[0]))
                    .toArray(Element[]::new);
        }
        return new Contraption(elements);
    }

    int rowCount() {
        return elements.length;
    }

    int columnCount() {
        return elements[0].length;
    }

    int maximalEnergizedPositionsCount() {
        final List<BeamPart> startPositions = new ArrayList<>();
        for (int column = 0; column < columnCount(); column++) {
            startPositions.add(new BeamPart(new Position(0, column), Direction.SOUTH));
        }
        for (int row = 0; row < rowCount(); row++) {
            startPositions.add(new BeamPart(new Position(row, 0), Direction.EAST));
        }
        for (int column = 0; column < columnCount(); column++) {
            startPositions.add(new BeamPart(new Position(rowCount() - 1, column), Direction.NORTH));
        }
        for (int row = 0; row < rowCount(); row++) {
            startPositions.add(new BeamPart(new Position(row, columnCount() - 1), Direction.WEST));
        }
        return startPositions.parallelStream()
                .mapToInt(this::energizedPositionsCount)
                .max()
                .orElse(0);
    }

    int energizedPositionsCount() {
        return energizedPositionsCount(new BeamPart(new Position(0, 0), Direction.EAST));
    }

    int energizedPositionsCount(final BeamPart start) {
        final var beam = new HashSet<BeamPart>();
        var currentBeamParts = List.of(start);
        while (!currentBeamParts.isEmpty()) {
            beam.addAll(currentBeamParts);
            currentBeamParts = currentBeamParts.stream()
                    .flatMap(BeamPart::nextParts)
                    .filter(not(beam::contains))
                    .toList();
        }
        return beam.size();
    }

    private Element elementAt(final Position position) {
        return elements[position.row()][position.column()];
    }

    private boolean contains(final Position position) {
        return position.row() >= 0 && position.row() < rowCount() &&
                position.column() >= 0 && position.column() < columnCount();
    }
}

void main() throws IOException {
    final var input = Files.readString(Path.of("input.txt"));
    final var contraption = Contraption.valueOf(input);

    final int energizedPositionsCount = contraption.energizedPositionsCount();
    System.out.println(energizedPositionsCount + " energized positions (part 1)");

    final long beforeInMs = System.currentTimeMillis();
    final int maximalEnergizedPositionsCount = contraption.maximalEnergizedPositionsCount();
    final long afterInMs = System.currentTimeMillis();
    System.out.println(maximalEnergizedPositionsCount + " energized positions maximum (part 2)");
    System.out.println("Part 2 execution time: " + (afterInMs - beforeInMs) + " ms");
}