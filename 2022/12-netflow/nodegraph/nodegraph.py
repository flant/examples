#!/usr/bin/env python

import json
import requests
from http.server import HTTPServer, BaseHTTPRequestHandler
from urllib.parse import urlparse, parse_qs
import time


class handler(BaseHTTPRequestHandler):
    def _set_headers(self):
        self.send_response(200)
        self.send_header("Content-type", "aplication/json")
        self.end_headers()

    def _fields(self):
        ret = '{"edges_fields": [{"field_name": "id","type": "string"},{"field_name": "source","type": "string"},{"field_name": "target","type": "string"}], "nodes_fields": [{"field_name": "id","type": "string"},{"field_name": "title","type": "string"},{"field_name": "mainStat","type": "number"},{"field_name": "secondaryStat","type": "number"}]}'
        return ret.encode("utf8") 
    def _data(self):
        s = set()
        e = []

        with open('/var/run/secrets/kubernetes.io/serviceaccount/token', 'r') as file:
            auth_token = file.read().replace('\n', '')
        auth_headers = {'Authorization': 'Bearer ' + auth_token}
        print(self.path)
        req = urlparse(self.path)
        args = parse_qs(req.query)
        filter = []
        if 'srcnamespace' in args:
          filter.append('srcnamespace="' + args['srcnamespace'][0] + '"')
        if 'dstnamespace' in args:
          filter.append('dstnamespace="' + args['dstnamespace'][0] + '"')

        t_now = int( time.time() )
        t_from = int(int(args['from'][0])/1000)
        t_to   = int(int(args['to'][0])/1000)
        if t_to < t_now:
          offset = ' offset ' + str(t_now - t_to) + 's'
        else:
          offset = ''
        
        window = t_to - t_from
        if window < 60:
            window = 60


        params = [ ( 'query', 'sum by (dsthost,srchost) (rate(netflow_connection_bytes{' + ','.join(filter) + '}[' + str(window) + 's]' + offset +'))' ) ]
        url = 'https://prometheus.d8-monitoring.svc.cluster.local:9090/api/v1/query'
        resp = requests.get(url, verify=False, headers=auth_headers, params=params )
        bytes_in = {} 
        bytes_out = {}
        data = resp.json()
        for i in data['data']['result']:
            tmp = {}
            tmp['source'] = i['metric']['srchost']
            tmp['target'] = i['metric']['dsthost']
            e.append(tmp)
            s.add(i['metric']['srchost'])
            s.add(i['metric']['dsthost'])
            
            if i['metric']['srchost'] not in bytes_out.keys():
                bytes_out[i['metric']['srchost']] = 0
            if i['metric']['dsthost'] not in bytes_in.keys():
                bytes_in[i['metric']['dsthost']] = 0 
            if i['metric']['dsthost'] not in bytes_out.keys():
                bytes_out[i['metric']['dsthost']] = 0
            if i['metric']['srchost'] not in bytes_in.keys():
                bytes_in[i['metric']['srchost']] = 0

            bytes_out[i['metric']['srchost']] = bytes_out[i['metric']['srchost']]+ float(i['value'][1])
            bytes_in[i['metric']['dsthost']] = bytes_in[i['metric']['dsthost']] + float(i['value'][1])
        x = []
        for d in s:
            tmp ={}
            tmp['title']=d
            tmp['id']=d
            tmp['mainStat']=round(bytes_in[d]*8/1024/10.24)/100
            tmp['secondaryStat']=round(bytes_out[d]*8/1024/10.24)/100            
            x.append(tmp)
        obj = {}
        obj['nodes'] = x
        obj['edges'] = e
        return json.dumps(obj).encode("utf8") 


    def do_GET(self):
        self._set_headers()
        req = urlparse(self.path)
        if (req.path == "/api/graph/fields"):
            self.wfile.write(self._fields())
        elif (req.path == "/api/graph/data"):
            self.wfile.write(self._data())
        else:
            self.wfile.write('{"result":"ok"}'.encode("utf8"))


if __name__ == "__main__":

    server_class=HTTPServer
    addr="0.0.0.0"
    appport=5000
    server_address = (addr, appport)
    httpd = server_class(server_address, handler)
    print(f"Starting httpd server on {addr}:{appport}")
    httpd.serve_forever()
