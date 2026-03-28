## ADDED Requirements

### Requirement: Date picker segmented input mode
The system SHALL provide a segmented input mode for editing `date` properties, consisting of three editable segments: year (YYYY), month (MM), and day (DD). The focused segment SHALL be visually highlighted.

#### Scenario: Activate segmented input on date property
- **WHEN** the user activates editing on a `date` property (Enter in properties panel or table cell)
- **THEN** the date picker SHALL open in segmented input mode with three segments: YYYY - MM - DD

#### Scenario: Pre-fill with current value
- **WHEN** the date property has an existing value
- **THEN** the segments SHALL be pre-filled with the existing date's year, month, and day

#### Scenario: Pre-fill with today when empty
- **WHEN** the date property has no value (empty/nil)
- **THEN** the segments SHALL be pre-filled with today's date

#### Scenario: Focus starts on year segment
- **WHEN** the date picker opens in segmented input mode
- **THEN** the year segment SHALL be focused initially

#### Scenario: Navigate segments with arrow keys
- **WHEN** the user presses left/right arrow or Tab/Shift+Tab
- **THEN** the focus SHALL move between year, month, and day segments

#### Scenario: Navigate segments wraps
- **WHEN** the focus is on the day segment and the user presses right arrow or Tab
- **THEN** the focus SHALL remain on the day segment (no wrap)

### Requirement: Segment increment and decrement
The system SHALL allow incrementing and decrementing segment values using up/down arrow keys, with automatic carry across segment boundaries.

#### Scenario: Increment year
- **WHEN** the year segment is focused and the user presses up arrow
- **THEN** the year value SHALL increase by 1

#### Scenario: Decrement year
- **WHEN** the year segment is focused and the user presses down arrow
- **THEN** the year value SHALL decrease by 1

#### Scenario: Increment month with carry
- **WHEN** the month segment is focused, the value is 12, and the user presses up arrow
- **THEN** the month SHALL become 1 and the year SHALL increase by 1

#### Scenario: Decrement month with carry
- **WHEN** the month segment is focused, the value is 1, and the user presses down arrow
- **THEN** the month SHALL become 12 and the year SHALL decrease by 1

#### Scenario: Increment day with carry
- **WHEN** the day segment is focused, the value is the last day of the month, and the user presses up arrow
- **THEN** the day SHALL become 1 and the month SHALL increase by 1 (with year carry if needed)

#### Scenario: Decrement day with carry
- **WHEN** the day segment is focused, the value is 1, and the user presses down arrow
- **THEN** the day SHALL become the last day of the previous month and the month SHALL decrease by 1 (with year carry if needed)

#### Scenario: Day clamping on month change
- **WHEN** the month changes and the current day exceeds the new month's last day
- **THEN** the day SHALL be clamped to the new month's last day (e.g., Jan 31 → Feb 28)

### Requirement: Direct digit input in segments
The system SHALL allow direct digit input to replace segment values.

#### Scenario: Type digits to replace segment
- **WHEN** the user types a digit while a segment is focused
- **THEN** the segment value SHALL be updated with the entered digits

#### Scenario: Two-digit month input
- **WHEN** the user types "03" in the month segment
- **THEN** the month SHALL become 3 and focus SHALL advance to the day segment

#### Scenario: Two-digit day input
- **WHEN** the user types "15" in the day segment
- **THEN** the day SHALL become 15

#### Scenario: Four-digit year input
- **WHEN** the user types "2025" in the year segment
- **THEN** the year SHALL become 2025 and focus SHALL advance to the month segment

### Requirement: Day-of-week feedback
The system SHALL display the day-of-week abbreviation as live feedback during segmented input.

#### Scenario: Show day of week
- **WHEN** the date picker is in segmented input mode
- **THEN** a day-of-week abbreviation (e.g., "Thu") SHALL be displayed to the right of the segments

#### Scenario: Day of week updates on change
- **WHEN** any segment value changes
- **THEN** the day-of-week abbreviation SHALL update immediately to reflect the new date

### Requirement: Date picker calendar mode
The system SHALL provide an inline calendar mode for date selection, displaying a mini month grid.

#### Scenario: Switch to calendar mode
- **WHEN** the user presses `c` in segmented input mode
- **THEN** the date picker SHALL switch to calendar mode displaying the current month

#### Scenario: Switch back to segmented input
- **WHEN** the user presses `c` in calendar mode
- **THEN** the date picker SHALL switch back to segmented input mode with the currently highlighted date

#### Scenario: Calendar grid layout
- **WHEN** the calendar is displayed
- **THEN** it SHALL show a 7-column grid with Mo Tu We Th Fr Sa Su header and up to 6 week rows

#### Scenario: Current date highlighted
- **WHEN** the calendar opens
- **THEN** the selected date (from segments) SHALL be highlighted in the calendar

#### Scenario: Today marker
- **WHEN** the calendar is displayed and today's date is visible in the current month
- **THEN** today's date SHALL be visually marked (distinct from the cursor highlight)

### Requirement: Calendar navigation
The system SHALL support keyboard navigation within the calendar grid.

#### Scenario: Navigate days with arrow keys
- **WHEN** the calendar is displayed and the user presses arrow keys or h/j/k/l
- **THEN** the highlight SHALL move one day in the corresponding direction

#### Scenario: Navigate across week boundary
- **WHEN** the highlight is on Sunday and the user presses right arrow
- **THEN** the highlight SHALL move to Monday of the next week

#### Scenario: Navigate across month boundary
- **WHEN** the highlight is on the last day of the month and the user presses right arrow or down arrow past the boundary
- **THEN** the calendar SHALL advance to the next month and highlight the first day

#### Scenario: Switch months with H/L
- **WHEN** the user presses `H` in calendar mode
- **THEN** the calendar SHALL show the previous month (maintaining the day or clamping)

#### Scenario: Switch months forward with L
- **WHEN** the user presses `L` in calendar mode
- **THEN** the calendar SHALL show the next month (maintaining the day or clamping)

#### Scenario: Jump to today with t
- **WHEN** the user presses `t` in calendar mode
- **THEN** the calendar SHALL navigate to today's month and highlight today's date

### Requirement: Date picker confirm and cancel
The system SHALL confirm the selected date on Enter and cancel on Esc, in both modes.

#### Scenario: Confirm in segmented input
- **WHEN** the user presses Enter in segmented input mode
- **THEN** the date SHALL be validated and saved as a YYYY-MM-DD string

#### Scenario: Confirm in calendar mode
- **WHEN** the user presses Enter in calendar mode
- **THEN** the highlighted date SHALL be validated and saved as a YYYY-MM-DD string

#### Scenario: Cancel in segmented input
- **WHEN** the user presses Esc in segmented input mode
- **THEN** the edit SHALL be cancelled and the original value restored

#### Scenario: Cancel in calendar mode
- **WHEN** the user presses Esc in calendar mode
- **THEN** the edit SHALL be cancelled and the original value restored

### Requirement: Date picker help bar integration
The help bar SHALL reflect the current date picker mode.

#### Scenario: Help bar in segmented input
- **WHEN** the date picker is in segmented input mode
- **THEN** the help bar SHALL show `[DATE]`

#### Scenario: Help bar in calendar mode
- **WHEN** the date picker is in calendar mode
- **THEN** the help bar SHALL show `[CAL]`

### Requirement: Leap year handling
The date picker SHALL correctly handle leap years in both modes.

#### Scenario: February 29 in leap year
- **WHEN** the year is a leap year and the month is February
- **THEN** day 29 SHALL be selectable in both segmented input and calendar

#### Scenario: February 29 in non-leap year
- **WHEN** the year is not a leap year and the month is February
- **THEN** the maximum day SHALL be 28 and day 29 SHALL NOT be selectable
- **AND** if the current day is 29, it SHALL be clamped to 28
