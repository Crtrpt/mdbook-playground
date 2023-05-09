import json
import sys

if __name__ == '__main__':
    # load both the context and the book representations from stdin
    book = json.load(sys.stdin)
    # context, book = json.load(sys.stdin)
    # # and now, we can just modify the content of the first chapter
  
    book['book']['sections'][0]['Chapter']['content'] = book['book']['sections'][0]['Chapter']['content']+"========"
    print(json.dumps(book,indent=2))
    # sys.stderr.write(json.dumps(book,indent=2))
    # # we are done with the book's modification, we can just print it to stdout, 
    
    #sys.stderr.write(json.dumps(book,indent=2))

