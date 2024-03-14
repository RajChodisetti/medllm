import requests
import io
from PyPDF2 import PdfReader

def extract_text_from_last_pages(pdf_url, output_path, num_pages=8):
    response = requests.get(pdf_url)
    pdf_file = io.BytesIO(response.content)

    reader = PdfReader(pdf_file)
    total_pages = len(reader.pages)
    text = ''
        
    for i in range(total_pages - num_pages, total_pages):
        page = reader.pages[i]
        text += page.extract_text()
    
    # Remove empty lines
    text_lines = filter(lambda x: x.strip(), text.split('\n'))
    text = '\n'.join(text_lines)
        
    with open(output_path, 'w') as txt_file:
        txt_file.write(text)

# Replace 'your_pdf_url' and 'output.txt' with your URL and desired file paths
pdf_url = "https://www.accessdata.fda.gov/drugsatfda_docs/label/2023/204311s000lbl.pdf#page=34"
output_path = 'metadata/output.txt'

extract_text_from_last_pages(pdf_url, output_path)
