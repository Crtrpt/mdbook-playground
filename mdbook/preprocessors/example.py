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
    context, book = json.load(sys.stdin)
    
    
    content=book['sections'][0]['Chapter']['content'] 
    x=re.sub(r'```kk(.*?)```',replace_func,content,0,re.MULTILINE|re.DOTALL)
    book['sections'][0]['Chapter']['content']=x

 
    sys.stderr.write(json.dumps(context,indent=2))
    print(json.dumps(book))
