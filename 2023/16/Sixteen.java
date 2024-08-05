import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Arrays;
import java.util.HashSet;
import java.util.Objects;
import java.util.Set;

import static java.util.function.Predicate.not;
import static java.util.stream.Collectors.toSet;

enum Direction {
    NORTH,
    SOUTH,
    EAST,
    WEST;
}

enum Element {
    EMPTY('.'),
    SPLITTER_NS('|'),
    SPLITTER_WE('-'),
    MIRROR_ES_NW('\\'),
    MIRROR_EN_SW('/');

    private static final Element[] CACHED_VALUES = values();

    private final char symbol;

    Element(final char symbol) {
        this.symbol = symbol;
    }

    static Element valueOf(final char symbol) {
        return Arrays.stream(CACHED_VALUES)
                .filter(element -> element.symbol == symbol)
                .findFirst()
                .orElseThrow();
    }

    char symbol() {
        return symbol;
    }

    Set<Direction> deviateBeamGoingTo(final Direction direction) {
        return switch (this) {
            case EMPTY -> Set.of(direction);
            case MIRROR_ES_NW -> switch (direction) {
                case NORTH -> Set.of(Direction.WEST);
                case SOUTH -> Set.of(Direction.EAST);
                case EAST -> Set.of(Direction.SOUTH);
                case WEST -> Set.of(Direction.NORTH);
            };
            case MIRROR_EN_SW -> switch (direction) {
                case NORTH -> Set.of(Direction.EAST);
                case SOUTH -> Set.of(Direction.WEST);
                case EAST -> Set.of(Direction.NORTH);
                case WEST -> Set.of(Direction.SOUTH);
            };
            case SPLITTER_NS -> switch (direction) {
                case EAST, WEST -> Set.of(Direction.NORTH, Direction.SOUTH);
                case NORTH, SOUTH -> Set.of(direction);
            };
            case SPLITTER_WE -> switch (direction) {
                case NORTH, SOUTH -> Set.of(Direction.EAST, Direction.WEST);
                case WEST, EAST -> Set.of(direction);
            };
        };
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

        Set<BeamPart> next() {
            return elementAt(position)
                    .deviateBeamGoingTo(direction).stream()
                    .map(newDirection -> new BeamPart(position.to(newDirection), newDirection))
                    .filter(beamPart -> Contraption.this.contains(beamPart.position))
                    .collect(toSet());
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

    Set<Position> energizedPositions() {
        final var beam = new HashSet<BeamPart>();

        var currentBeamParts = Set.of(new BeamPart(new Position(0, 0), Direction.EAST));
        while (!currentBeamParts.isEmpty()) {
            beam.addAll(currentBeamParts);
            currentBeamParts = currentBeamParts.stream()
                    .flatMap(beamPart -> beamPart.next().stream())
                    .filter(not(beam::contains))
                    .collect(toSet());
        }

        return beam.stream().map(beamPart -> beamPart.position).collect(toSet());
    }

    private Element elementAt(final Position position) {
        return elements[position.row()][position.column()];
    }

    private boolean contains(final Position position) {
        return position.row() >= 0 && position.row() < rowCount() &&
                position.column() >= 0 && position.column() < columnCount();
    }

    @Override
    public String toString() {
        final var sb = new StringBuilder();
        for (final var row : elements) {
            for (final var element : row) {
                sb.append(element.symbol());
            }
            sb.append('\n');
        }
        return sb.toString();
    }
}

void main() throws IOException {
    final var input = Files.readString(Path.of("input.txt"));
    final var contraption = Contraption.valueOf(input);
    final Set<Position> energizedPositions = contraption.energizedPositions();
    System.out.println(energizedPositions.size() + " energized positions");
}