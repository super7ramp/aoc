import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.*;

import static java.util.stream.Collectors.toMap;

enum Direction {
    NORTH,
    EAST,
    SOUTH,
    WEST,
}

enum Tile {

    VERTICAL('|', Direction.NORTH, Direction.SOUTH),
    HORIZONTAL('-', Direction.EAST, Direction.WEST),
    BEND_NORTH_EAST('L', Direction.NORTH, Direction.EAST),
    BEND_NORTH_WEST('J', Direction.NORTH, Direction.WEST),
    BEND_SOUTH_WEST('7', Direction.SOUTH, Direction.WEST),
    BEND_SOUTH_EAST('F', Direction.SOUTH, Direction.EAST),
    GROUND('.'),
    START('S', Direction.values());

    private static final Tile[] CACHED_VALUES = values();
    private final char symbol;
    private final List<Direction> directions;

    Tile(final char symbol, final Direction... directions) {
        this.symbol = symbol;
        this.directions = List.of(directions);
    }

    static Tile valueOf(final char symbol) {
        return Arrays.stream(CACHED_VALUES)
                .filter(tile -> tile.symbol == symbol)
                .findFirst()
                .orElseThrow();
    }

    static Tile fromDirections(final Direction... directions) {
        return Arrays.stream(CACHED_VALUES)
                .filter(tile -> new HashSet<>(tile.directions).equals(Set.of(directions)))
                .findFirst()
                .orElseThrow();
    }

    List<Direction> directions() {
        return directions;
    }

    boolean isStart() {
        return this == START;
    }

    boolean isHorizontal() {
        return this == HORIZONTAL;
    }

    boolean isBendNorthWest() {
        return this == BEND_NORTH_WEST;
    }

    boolean isBendSouthWest() {
        return this == BEND_SOUTH_WEST;
    }
}


record Position(int row, int column) {

    Position neighborPositionTo(final Direction direction) {
        return switch (direction) {
            case NORTH -> new Position(row - 1, column);
            case EAST -> new Position(row, column + 1);
            case SOUTH -> new Position(row + 1, column);
            case WEST -> new Position(row, column - 1);
        };
    }

    Direction directionTo(final Position other) {
        if (other.row == row + 1 && other.column == column) {
            return Direction.SOUTH;
        }
        if (other.row == row - 1 && other.column == column) {
            return Direction.NORTH;
        }
        if (other.row == row && other.column == column + 1) {
            return Direction.EAST;
        }
        if (other.row == row && other.column == column - 1) {
            return Direction.WEST;
        }
        throw new IllegalArgumentException("Not a neighbor position");
    }
}

static class Polygon {

    private final Map<Position, Tile> edges;

    Polygon(final Map<Position, Tile> edges) {
        this.edges = edges;
    }

    boolean encloses(final Position position) {
        if (edges.containsKey(position)) {
            return false;
        }
        return countIntersectionsWithRayCast(position) % 2 == 1;
    }

    // See https://en.wikipedia.org/wiki/Point_in_polygon.
    private int countIntersectionsWithRayCast(final Position position) {
        int count = 0;
        Tile previous = null;
        for (Position current = new Position(position.row(), 0);
             !current.equals(position);
             current = new Position(position.row(), current.column() + 1)) {
            final Tile edge = edges.get(current);
            if (edge != null && !edge.isHorizontal()) {
                count++;
                // Special cases: JF and 7L
                if (edge.isBendNorthWest() && previous == Tile.BEND_SOUTH_EAST ||
                        edge.isBendSouthWest() && previous == Tile.BEND_NORTH_EAST) {
                    count++;
                }
                previous = edge;
            }
        }
        return count;
    }
}

static class PipeMap {

    private final Tile[][] map;

    private PipeMap(final Tile[][] map) {
        this.map = map;
    }

    static PipeMap valueOf(final String input) {
        final var lines = input.split("\n");
        final var map = new Tile[lines.length][];
        for (int rowIndex = 0; rowIndex < lines.length; rowIndex++) {
            final var line = lines[rowIndex];
            final var row = new Tile[line.length()];
            for (int columnIndex = 0; columnIndex < line.length(); columnIndex++) {
                row[columnIndex] = Tile.valueOf(line.charAt(columnIndex));
            }
            map[rowIndex] = row;
        }
        return new PipeMap(map);
    }

    List<Position> findLoop() {
        final var loop = new ArrayList<Position>();
        final Position start = findStartPosition();
        for (Position current = start, previous = null; !current.equals(start) || loop.isEmpty(); ) {
            final Position oldCurrent = current;
            current = findNextPosition(current, previous);
            previous = oldCurrent;
            loop.add(current);
        }
        return loop;
    }

    private Position findStartPosition() {
        for (int rowIndex = 0; rowIndex < map.length; rowIndex++) {
            for (int columnIndex = 0; columnIndex < map[rowIndex].length; columnIndex++) {
                if (map[rowIndex][columnIndex].isStart()) {
                    return new Position(rowIndex, columnIndex);
                }
            }
        }
        throw new IllegalStateException("No start position found");
    }

    private Position findNextPosition(final Position currentPosition, final Position previousPosition) {
        return tileAt(currentPosition)
                .directions()
                .stream()
                .map(currentPosition::neighborPositionTo)
                .filter(this::contains)
                .filter(neighborPosition -> tileAt(neighborPosition)
                        .directions()
                        .contains(neighborPosition.directionTo(currentPosition)))
                .filter(neighborPosition -> !neighborPosition.equals(previousPosition))
                .findFirst()
                .orElseThrow();
    }

    private Tile tileAt(final Position position) {
        return map[position.row()][position.column()];
    }

    private boolean contains(final Position position) {
        return position.row() >= 0 && position.row() < map.length
                && position.column() >= 0 && position.column() < map[position.row()].length;
    }

    List<Position> enclosedTilePositions(final List<Position> loopPositions) {
        final Map<Position, Tile> tilesWithPositions = loopPositions.stream()
                .collect(toMap(position -> position, this::tileAt));

        // Replace start tile with the tile hidden behind it.
        final Position startPosition = loopPositions.getLast();
        final Position secondPosition = loopPositions.getFirst();
        final Position lastPosition = loopPositions.get(loopPositions.size() - 2);
        final Tile actualStartTile = Tile.fromDirections(startPosition.directionTo(secondPosition),
                startPosition.directionTo(lastPosition));
        tilesWithPositions.put(startPosition, actualStartTile);

        final var loop = new Polygon(tilesWithPositions);

        final var enclosedTilePositions = new ArrayList<Position>();
        for (int row = 0; row < map.length; row++) {
            for (int column = 0; column < map[row].length; column++) {
                final Position position = new Position(row, column);
                if (loop.encloses(position)) {
                    enclosedTilePositions.add(position);
                }
            }
        }

        return enclosedTilePositions;
    }
}

void main() throws IOException {
    final String input = Files.readString(Path.of("input.txt"));
    final PipeMap pipeMap = PipeMap.valueOf(input);

    final List<Position> loopPositions = pipeMap.findLoop();
    System.out.println("Loop positions: " + loopPositions);
    System.out.println("Loop position count: " + loopPositions.size());
    System.out.println("Farthest position: " + loopPositions.size() / 2);

    final List<Position> enclosedTilePositions = pipeMap.enclosedTilePositions(loopPositions);
    System.out.println("Enclosed tile positions: " + enclosedTilePositions);
    System.out.println("Enclosed tile count: " + enclosedTilePositions.size());
}