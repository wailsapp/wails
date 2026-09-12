# Linux/macOS orchestration regression: python3 check.py /absolute/path/to/wails3
# Uses lightweight task stand-ins to check launch ordering without a GUI/npm.
import os,pathlib,signal,socket,subprocess,sys,tempfile,time
timeout_case = len(sys.argv) > 2 and sys.argv[2] == 'timeout'
with tempfile.TemporaryDirectory(prefix='wails-readiness-') as d:
 root=pathlib.Path(d)
 with socket.socket() as s:s.bind(('127.0.0.1',0));port=s.getsockname()[1]
 fake=root/'wails3'
 fake.write_text('''#!/usr/bin/env python3
import os,socket,sys,time,pathlib
if sys.argv[1:]==['task','common:dev:frontend']:
 time.sleep(6)
 s=socket.socket();s.bind(('127.0.0.1',int(os.environ['WAILS_VITE_PORT'])));s.listen()
 while True:
  c,_=s.accept();c.close()
elif sys.argv[1:]==['task','run']:
 try:
  with socket.create_connection(('localhost',int(os.environ['WAILS_VITE_PORT'])),timeout=1):pass
  pathlib.Path('result').write_text('READY')
 except OSError:pathlib.Path('result').write_text('STARTED BEFORE FRONTEND')
 while True:time.sleep(1)
''');fake.chmod(0o755)
 (root/'config.yml').write_text('''dev_mode:
  root_path: .
  log_level: error
  ignore:
    watched_extension: ["*.go"]
  executes:
    - cmd: wails3 build DEV=true
      type: blocking
    - cmd: wails3 task common:dev:frontend
      type: background
    - cmd: wails3 task run
      type: primary
''')
 if timeout_case:
  config=root/'config.yml'
  config.write_text(config.read_text().replace('      type: background', '      type: background\n      readiness:\n        tcp: localhost:'+str(port)+'\n        timeout: 200ms'))
 with open(root/'log','w') as log:
  p=subprocess.Popen([sys.argv[1],'dev','--config',str(root/'config.yml'),'--port',str(port)],cwd=root,env={**os.environ,'PATH':str(root)+':'+os.environ['PATH']},stdout=log,stderr=log,start_new_session=True)
  try:
   deadline=time.monotonic()+15
   while not (root/'result').exists() and p.poll() is None and time.monotonic()<deadline:time.sleep(.05)
   result=(root/'result').read_text() if (root/'result').exists() else 'NO APP RESULT'
   exit_code=p.poll()
  finally:
   if p.poll() is None:os.killpg(p.pid,signal.SIGTERM)
   try:p.wait(timeout=3)
   except subprocess.TimeoutExpired:os.killpg(p.pid,signal.SIGKILL);p.wait()
 if timeout_case:
  if result=='NO APP RESULT' and exit_code not in (None,0) and 'did not become ready' in (root/'log').read_text():
   print('TIMEOUT: app correctly not started');sys.exit(0)
  print(result);print((root/'log').read_text());sys.exit(1)
 print(result)
 if result!='READY':print((root/'log').read_text());sys.exit(1)
