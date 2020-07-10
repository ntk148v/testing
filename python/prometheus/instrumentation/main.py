import prometheus_client as prom
import random
import time

req_summary = prom.Summary('python_my_req_example',
                           'Time spent processing a request')


@req_summary.time()
def process_request(t):
    time.sleep(t)


if __name__ == '__main__':
    counter = prom.Counter('python_my_counter', 'This is my counter')
    gauge = prom.Gauge('python_my_gauge', 'This is my gauge')
    histogram = prom.Histogram('python_my_histogram', 'This is my histogram')
    summary = prom.Summary('python_my_summary', 'This is my summary')
    prom.start_http_server(8080)

    while True:
        counter.inc(random.random())
        gauge.set(random.random() * 15 - 5)
        histogram.observe(random.random() * 10)
        summary.observe(random.random() * 10)
        process_request(random.random() * 5)

        time.sleep(1)
