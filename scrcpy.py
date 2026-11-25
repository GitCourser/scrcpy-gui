# ==============================================================================
# Author       : Courser
# Date         : 2023-03-16 21:54:42
# LastEditTime : 2025-04-11 10:28:08
# Description  : scrcpy GUI 0.6
#                0.1 tkinter
#                0.2 PySide2
#                0.3 添加拖放路径, 添加扫描进度, 优化扫描
#                0.4 添加 scrcpy v2.1 新选项 "仅录音"
#                0.5 只允许启动一个实例
#                0.6 添加 scrcpy v2.5 新选项 "鼠标悬停"
#                0.7 更新为 Pyside6
# ==============================================================================

import json
from subprocess import CREATE_NO_WINDOW, PIPE, Popen, getoutput, run
from threading import Thread
from time import localtime, sleep, strftime

from pyperclip import copy
from PySide6.QtCore import QDir, QFileInfo, Qt, QThread, Signal, Slot
from PySide6.QtWidgets import (
    QApplication,
    QDialog,
    QFileDialog,
    QMessageBox,
    QProgressDialog,
)
from qdarkstyle import load_stylesheet_pyside2

from main_ui import Ui_main
from ps import scan_range, scan_single

version = 'scrcpy GUI v0.7'


def new_thread(func):
    """新线程装饰器"""

    def inner(*args, **kwargs):
        t = Thread(target=func, args=args, kwargs=kwargs)
        t.start()

    return inner


class ScanThread(QThread):
    # 定义信号：扫描完成时传递结果和耗时
    scan_finished = Signal(list, str)
    scan_failed = Signal(str)

    def __init__(self, ip, port_range, threads, timeout, port_filter):
        super().__init__()
        self.ip = ip
        self.port_range = port_range
        self.threads = threads
        self.timeout = timeout
        self.port_filter = port_filter

    def run(self):
        """线程主逻辑"""
        try:
            port_open, rs_time = scan_range(
                self.ip,
                self.port_range,
                self.threads,
                self.timeout,
                *self.port_filter,
            )
            self.scan_finished.emit(port_open, rs_time)
        except RuntimeError as e:
            self.scan_failed.emit(f"错误: {str(e)}")
        except Exception as e:
            self.scan_failed.emit(f"未知错误: {str(e)}")


class MainWindow(QDialog, Ui_main):
    # scan_finished = Signal()

    def __init__(self):
        super().__init__()
        self.setupUi(self)
        # 设置窗口标题
        self.setWindowTitle(version)
        # 只有关闭按钮
        self.setWindowFlags(Qt.WindowCloseButtonHint)
        # 晚上自动切换到暗色主题
        if localtime().tm_hour > 18 or localtime().tm_hour < 8:
            self.setStyleSheet(load_stylesheet_pyside2())
        # 读取配置
        self.cfg = self._get_cfg()
        # 初始化界面检查是否扫描
        self._init_ui()

    def dragEnterEvent(self, event):
        """拖入"""
        if event.mimeData().hasUrls():
            event.accept()
        else:
            event.ignore()

    def dropEvent(self, event):
        """取得拖入路径"""
        path = event.mimeData().text()[8:]
        if QFileInfo(path).isDir():
            path = QDir(path).path()
        if QFileInfo(path).isFile():
            path = QFileInfo(path).path()
        self.strPath.setText(path)

    def _init_ui(self):
        """初始化界面"""
        # 检查是否开启扫描
        if not self.cfg.get('isScan'):
            return
        # 先单独扫描上次使用的端口
        if scan_single(self.cfg['strIP'], int(self.cfg['strPort'])):
            self.setWindowTitle(f'{version} - 端口正常')
            return
        # 开启范围扫描并显示扫描进度
        self._scan()

    def _scan(self):
        self.pd = QProgressDialog('正在扫描...', None, 0, 0, self, Qt.WindowTitleHint)  # FramelessWindowHint
        self.pd.setStyleSheet('QProgressBar {text-align: center;}')
        self.pd.setWindowModality(Qt.WindowModal)
        self.pd.setWindowTitle('请等待')
        self.pd.setFixedSize(240, 80)
        # self.scan_finished.connect(self.pd.close)
        self.pd.show()
        # 创建并启动线程
        port_filter = self.strUnPort.text().split(',')
        self.scan_thread = ScanThread(
            self.cfg['strIP'],
            self.cfg['strPortRange'],
            self.numThread.value(),
            self.numTimeout.value(),
            port_filter
        )
        self.scan_thread.scan_finished.connect(self.on_scan_finished)
        self.scan_thread.scan_failed.connect(self.on_scan_failed)
        self.scan_thread.start()

    def on_scan_finished(self, port_open, rs_time):
        """扫描成功槽函数"""
        if len(port_open) == 1:
            self.strPort.setText(str(port_open[0]))
            self.setWindowTitle(f'{version} - 扫描成功 {rs_time}')
        else:
            self.setWindowTitle(f'{version} - 扫描结果 {port_open} {rs_time}')
            copy(','.join(map(str, port_open)))
        self.pd.close()

    def on_scan_failed(self, error_msg):
        """扫描失败槽函数"""
        self.setWindowTitle(f'{version} - {error_msg}')
        self.pd.close()
        QMessageBox.critical(self, '错误', error_msg)

    # @new_thread
    # def _scan(self):
    #     """扫描端口"""
    #     port_filter = self.strUnPort.text().split(',')
    #     print('过滤端口', port_filter)
    #     try:
    #         print('开始扫描...')
    #         port_open, rs_time = scan_range(
    #             self.cfg['strIP'],
    #             self.cfg['strPortRange'],
    #             self.numThread.value(),
    #             self.numTimeout.value(),
    #             *port_filter,
    #         )
    #     except RuntimeError:
    #         print('错误: 线程过多')
    #         self.setWindowTitle(f'{version} - 错误: 线程过多, 请减少线程数')
    #     else:
    #         print('扫描完毕')
    #         if len(port_open) == 1:
    #             self.strPort.setText(str(port_open[0]))
    #             self.setWindowTitle(f'{version} - 扫描成功 {rs_time}')
    #         else:
    #             self.setWindowTitle(f'{version} - 扫描结果 {port_open} {rs_time}')
    #             copy(','.join(map(str, port_open)))
    #     finally:
    #         self.scan_finished.emit()

    def _get_cfg(self):
        """读取配置"""
        try:
            with open('scrcpy.cfg', 'r') as f:
                cfg = json.load(f)
        except FileNotFoundError:
            return {}
        for i, j in cfg.items():
            try:
                if i[:3] == 'str':
                    eval(f'self.{i}.setText(j)')
                elif i[:2] == 'is':
                    eval(f'self.{i}.setChecked(j)')
                elif i[:3] == 'num':
                    eval(f'self.{i}.setValue(j)')
            except AttributeError:
                continue
        return cfg

    def _set_cfg(self):
        """保存配置"""
        cfglist = [
            'strPath',
            'strIP',
            'strPort',
            'isScan',
            'strPortRange',
            'strUnPort',
            'numThread',
            'numTimeout',
            'isTop',
            'isSw',
            'isPoweroff',
            'isHEVC',
            'isNoAudio',
            'isNoMousehover',
            'isRecord',
            'isRecordAudio',
            'strKey',
            'strSize',
            'strBitrate',
            'strBuffer',
        ]
        for i in cfglist:
            if i[:3] == 'str':
                self.cfg[i] = eval(f'self.{i}.text()')
            elif i[:2] == 'is':
                self.cfg[i] = eval(f'self.{i}.isChecked()')
            elif i[:3] == 'num':
                self.cfg[i] = eval(f'self.{i}.value()')
            if any([self.cfg[i] == x for x in (False, '', '-', 400, 1.5)]):
                del self.cfg[i]
        with open('scrcpy.cfg', 'w') as f:
            json.dump(self.cfg, f, indent=2)

    def _create_cmd(self, path):
        """生成命令"""
        self.ip = f'''{self.cfg['strIP']}:{self.cfg['strPort']}'''
        self.cmd_dev = [f'{path}adb', 'devices']
        self.cmd_conn = [f'{path}adb', 'connect', self.ip]
        self.cmd_pair = [f'{path}adb', 'pair']
        self.cmd_scrcpy = [f'{path}scrcpy', f'-s{self.ip}', '--max-fps=60']
        self.isOnlyAudio = False
        if 'isTop' in self.cfg:
            self.cmd_scrcpy.append('--always-on-top')
        if 'isSw' in self.cfg:
            self.cmd_scrcpy.append('-Sw')
        if 'isPoweroff' in self.cfg:
            self.cmd_scrcpy.append('--power-off-on-close')
        if 'isHEVC' in self.cfg:
            self.cmd_scrcpy.append('--video-codec=h265')
        if 'isNoAudio' in self.cfg:
            self.cmd_scrcpy.append('--no-audio')
        if 'isNoMousehover' in self.cfg:
            self.cmd_scrcpy.append('--no-mouse-hover')
        if 'isRecord' in self.cfg:
            if not QDir().exists('rec'):
                QDir().mkdir('rec')
            self.cmd_scrcpy.append(strftime('--record=rec/%Y-%m-%d_%H-%M-%S.mp4'))
        if 'isRecordAudio' in self.cfg:
            if not QDir().exists('rec'):
                QDir().mkdir('rec')
            self.isOnlyAudio = True
            self.cmd_scrcpy.append('--no-video')
            self.cmd_scrcpy.append('--audio-buffer=100')
            # self.cmd_scrcpy.append('--audio-codec=aac')
            self.cmd_scrcpy.append(strftime('--record=rec/%Y-%m-%d_%H-%M-%S.opus'))
            try:
                self.cmd_scrcpy.remove('-Sw')
            except ValueError:
                pass
            try:
                self.cmd_scrcpy.remove('--power-off-on-close')
            except ValueError:
                pass
        if 'strKey' in self.cfg:
            self.cmd_scrcpy.append(f'''--shortcut-mod={self.cfg['strKey']}''')
        if 'strSize' in self.cfg:
            self.cmd_scrcpy.append(f'''-m{self.cfg['strSize']}''')
        if 'strBitrate' in self.cfg:
            self.cmd_scrcpy.append(f'''-b{self.cfg['strBitrate']}M''')
        if 'strBuffer' in self.cfg:
            self.cmd_scrcpy.append(f'''--display-buffer={self.cfg['strBuffer']}''')
        if self.isPair.isChecked():
            self.cmd_pair.append(f'{self.strIP.text()}:{self.strPairPort.text()}')
        print(self.cmd_scrcpy)

    @new_thread
    def _start(self):
        """启动 scrcpy"""

        # adb pair
        if self.isPair.isChecked():
            p = Popen(self.cmd_pair, stdin=PIPE, creationflags=CREATE_NO_WINDOW)
            p.stdin.write(self.strPairCode.text().encode())
            p.stdin.close()
            sleep(1)

        # adb devices
        if (
            f'''{self.ip}\tdevice'''
            not in run(self.cmd_dev, stdout=PIPE, creationflags=CREATE_NO_WINDOW).stdout.decode()
        ):
            # adb connect
            run(self.cmd_conn, creationflags=CREATE_NO_WINDOW)

        # scrcpy
        if not self.isOnlyAudio:
            run(self.cmd_scrcpy, stdout=PIPE, creationflags=CREATE_NO_WINDOW)
        else:
            print('开始录音, Ctrl+C 结束')
            run(self.cmd_scrcpy)

        # adb disconnect
        # run(['adb', 'disconnect'], creationflags=CREATE_NO_WINDOW)

    @Slot()
    def on_btnPath_clicked(self):
        """选择目录"""
        filePath = QFileDialog.getExistingDirectory(self)  # options=QFileDialog.DontUseNativeDialog
        self.strPath.setText(filePath)

    @Slot()
    def on_btnStart_clicked(self):
        """启动按钮"""

        # 检查
        if self.strIP.text() == '...' or not self.strPort.text():
            QMessageBox.warning(self, '警告', '请先设置IP和端口')
            return
        cmdpath = f'{self.strPath.text()}/' if self.strPath.text() else ''
        if not QFileInfo.exists(f'{cmdpath}scrcpy.exe'):
            QMessageBox.warning(self, '警告', '请先设置scrcpy目录')
            return
        if self.isScan.isChecked() and (self.strPortRange.text() == '-'):
            QMessageBox.warning(self, '警告', '开了扫描就必须设置端口范围')
            return
        if self.isPair.isChecked() and (not self.strPairPort.text() or not self.strPairCode.text()):
            QMessageBox.warning(self, '警告', '开了配对就必须设置配对端口和配对码')
            return

        # 保存参数并生成命令
        self._set_cfg()
        self._create_cmd(cmdpath)

        # 运行
        self._start()
        self.close()


def check_single_instance(process_name):
    """检查是否已经有实例在运行"""
    output = getoutput('tasklist')
    instances = output.count(process_name)
    if instances > 0:
        return True


if __name__ == '__main__':
    # QApplication.setAttribute(Qt.AA_EnableHighDpiScaling)
    app = QApplication([])
    if check_single_instance('scrcpy.exe'):
        print('不要重复启动程序！')
        QMessageBox.information(None, '提示', '不要重复启动程序！')
        exit()
    window = MainWindow()
    window.show()
    exit(app.exec())
