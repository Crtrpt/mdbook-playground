import json
import sys
import re
import base64

def replace_func(matchobj):
       content=matchobj.group(0)
       jsonStr=content[6:-3].rstrip().rstrip()
       base64Str=base64.b64encode(jsonStr.encode())
       return '<div class="kk" data-source='+base64Str.decode("ascii")+'></div>'

if __name__ == '__main__':
    # print("================================")
    # if len(sys.argv) > 1: # we check if we received any argument
    #     if sys.argv[1] == "supports": 
    #         # then we are good to return an exit status code of 0, since the other argument will just be the renderer's name
    #         sys.exit(0)

    # load both the context and the book representations from stdin
    context, book = json.load(sys.stdin)
    # and now, we can just modify the content of the first chapter
    
    content=book['sections'][0]['Chapter']['content'] 
    x=re.sub(r'```kk(.*?)```',replace_func,content,0,re.MULTILINE|re.DOTALL)
    book['sections'][0]['Chapter']['content']=x

 
    sys.stderr.write(json.dumps(context,indent=2))
    #sys.stderr.write(context.config.preprocessor.example.aa)
    # sys.stderr.write(json.dumps(context['config']['preprocessor']['example']['aa'],indent=2))
    # we are done with the book's modification, we can just print it to stdout, 
    print(json.dumps(book))
