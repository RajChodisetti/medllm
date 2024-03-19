import csv

def get_seventh_column(csv_file):
    with open(csv_file, 'r') as file:
        reader = csv.reader(file)
        # Skip header row if present
        next(reader, None)
        
        # Iterate over the first 10 rows
        rows = []
        for _ in range(10):
            try:
                row = next(reader)
                # Get the 7th column value and append to the list
                rows.append(row[6])  # Indexing is 0-based, so 6 represents the 7th column
            except StopIteration:
                break
        
        return rows

# Example usage:
csv_file = 'top10.csv'  # Replace with the path to your CSV file
seventh_column_values = get_seventh_column(csv_file)
print(seventh_column_values)
