# ==============================================================================
# Author       : Courser
# Date         : 2022-06-02 23:30:18
# LastEditTime : 2023-04-07 10:36:15
# Description  : 扫描无线调试端口
# ==============================================================================

import socket
from time import perf_counter
from queue import Queue
from threading import Thread

port_open = []


class PortScaner(object):
    class PortScan(Thread):
        def __init__(self, port_queue, ip, timeout=1):
            """初始化参数"""
            Thread.__init__(self)
            self.__port_queue = port_queue
            self.__ip = ip
            self.__timeout = timeout

        def run(self):
            ip = self.__ip
            timeout = self.__timeout
            # 如果端口队列不为空，循环执行
            while True:
                if self.__port_queue.empty():
                    break
                port = self.__port_queue.get()

                s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
                s.settimeout(timeout)
                result_code = s.connect_ex((ip, port))
                # 开放返回0
                if result_code == 0:
                    port_open.append(port)
                s.close()

    def get_port_lists(self, start_port=1, end_port=65535):
        """获取扫描的端口list，端口号无效返回[1-65535]"""
        if start_port >= 1 and end_port <= 65535 and start_port <= end_port:
            return list(range(start_port, end_port + 1))
        else:
            return list(range(1, 65535 + 1))


def scan_range(ip: str, port_range: str, thread_num: int, timeout, *args):
    """范围扫描

    Args:
        ip (str): ip地址
        port_range (str): 端口范围
        thread_num (int): 线程数
        timeout (Number): 超时

    Returns:
        (list, str): 端口列表, 扫描时间
    """
    start_time = perf_counter()
    port_scner = PortScaner()
    port_queue = Queue()
    threads = []

    # 获取端口列表
    start_port, end_port = port_range.split('-')
    port_list = port_scner.get_port_lists(int(start_port), int(end_port))
    for i in args:
        try:
            port_list.remove(int(i))
        except ValueError:
            pass

    # 将端口列表放入队列
    for port in port_list:
        port_queue.put(port)

    # 创建线程
    for _ in range(thread_num):
        threads.append(port_scner.PortScan(port_queue, ip, timeout))

    # 启动线程
    for thread in threads:
        thread.start()
    for thread in threads:
        thread.join()

    end_time = perf_counter()
    result_time = f'⏱️{end_time - start_time:.3f}s'
    print(port_open, result_time)
    return port_open, result_time


def scan_single(ip: str, port: int):
    """单个扫描

    Args:
        ip (str): ip地址
        port (int): 端口号

    Returns:
        str: 扫描时间
    """
    start_time = perf_counter()
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        result_code = s.connect_ex((ip, port))
        if result_code == 0:
            end_time = perf_counter()
            result_time = f'⏱️{end_time - start_time:.3f}s'
            print('open', result_time)
            return result_time


if __name__ == '__main__':
    from sys import argv
    if len(argv) == 1:
        print('单个端口扫描: ps 192.168.x.x 5555\n范围端口扫描: ps 192.168.x.x 35000-44000 线程数 [可选超时]')
    elif len(argv) == 3:
        scan_single(argv[1], int(argv[2]))
    else:
        scan_range(argv[1], argv[2], int(argv[3]), argv[4] if len(argv) == 5 else 1)
