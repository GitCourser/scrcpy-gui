# -*- coding: utf-8 -*-

################################################################################
## Form generated from reading UI file 'main.ui'
##
## Created by: Qt User Interface Compiler version 6.8.2
##
## WARNING! All changes made in this file will be lost when recompiling UI file!
################################################################################

from PySide6.QtCore import (QCoreApplication, QDate, QDateTime, QLocale,
    QMetaObject, QObject, QPoint, QRect,
    QSize, QTime, QUrl, Qt)
from PySide6.QtGui import (QBrush, QColor, QConicalGradient, QCursor,
    QFont, QFontDatabase, QGradient, QIcon,
    QImage, QKeySequence, QLinearGradient, QPainter,
    QPalette, QPixmap, QRadialGradient, QTransform)
from PySide6.QtWidgets import (QApplication, QCheckBox, QDialog, QDoubleSpinBox,
    QFrame, QGridLayout, QGroupBox, QHBoxLayout,
    QLabel, QLineEdit, QPushButton, QSizePolicy,
    QSpinBox, QVBoxLayout, QWidget)

class Ui_main(object):
    def setupUi(self, main):
        if not main.objectName():
            main.setObjectName(u"main")
        main.resize(600, 440)
        main.setMinimumSize(QSize(600, 440))
        main.setMaximumSize(QSize(600, 440))
        font = QFont()
        font.setFamilies([u"Microsoft YaHei UI"])
        main.setFont(font)
        main.setAcceptDrops(True)
        self.gridLayout_2 = QGridLayout(main)
        self.gridLayout_2.setObjectName(u"gridLayout_2")
        self.horizontalLayout = QHBoxLayout()
        self.horizontalLayout.setObjectName(u"horizontalLayout")
        self.label = QLabel(main)
        self.label.setObjectName(u"label")

        self.horizontalLayout.addWidget(self.label)

        self.strPath = QLineEdit(main)
        self.strPath.setObjectName(u"strPath")
        self.strPath.setReadOnly(True)

        self.horizontalLayout.addWidget(self.strPath)

        self.btnPath = QPushButton(main)
        self.btnPath.setObjectName(u"btnPath")
        self.btnPath.setAutoDefault(False)

        self.horizontalLayout.addWidget(self.btnPath)


        self.gridLayout_2.addLayout(self.horizontalLayout, 0, 0, 1, 1)

        self.line = QFrame(main)
        self.line.setObjectName(u"line")
        self.line.setFrameShape(QFrame.Shape.HLine)
        self.line.setFrameShadow(QFrame.Shadow.Sunken)

        self.gridLayout_2.addWidget(self.line, 1, 0, 1, 1)

        self.horizontalLayout_2 = QHBoxLayout()
        self.horizontalLayout_2.setObjectName(u"horizontalLayout_2")
        self.groupBox = QGroupBox(main)
        self.groupBox.setObjectName(u"groupBox")
        self.groupBox.setMinimumSize(QSize(188, 0))
        self.groupBox.setMaximumSize(QSize(192, 16777215))
        self.verticalLayout = QVBoxLayout(self.groupBox)
        self.verticalLayout.setObjectName(u"verticalLayout")
        self.label_2 = QLabel(self.groupBox)
        self.label_2.setObjectName(u"label_2")

        self.verticalLayout.addWidget(self.label_2)

        self.strIP = QLineEdit(self.groupBox)
        self.strIP.setObjectName(u"strIP")

        self.verticalLayout.addWidget(self.strIP)

        self.label_3 = QLabel(self.groupBox)
        self.label_3.setObjectName(u"label_3")

        self.verticalLayout.addWidget(self.label_3)

        self.strPort = QLineEdit(self.groupBox)
        self.strPort.setObjectName(u"strPort")

        self.verticalLayout.addWidget(self.strPort)


        self.horizontalLayout_2.addWidget(self.groupBox)

        self.isScan = QGroupBox(main)
        self.isScan.setObjectName(u"isScan")
        self.isScan.setMinimumSize(QSize(188, 0))
        self.isScan.setMaximumSize(QSize(192, 16777215))
        self.isScan.setCheckable(True)
        self.isScan.setChecked(False)
        self.gridLayout_4 = QGridLayout(self.isScan)
        self.gridLayout_4.setObjectName(u"gridLayout_4")
        self.label_13 = QLabel(self.isScan)
        self.label_13.setObjectName(u"label_13")

        self.gridLayout_4.addWidget(self.label_13, 4, 0, 1, 1)

        self.strPortRange = QLineEdit(self.isScan)
        self.strPortRange.setObjectName(u"strPortRange")

        self.gridLayout_4.addWidget(self.strPortRange, 1, 0, 1, 2)

        self.numTimeout = QDoubleSpinBox(self.isScan)
        self.numTimeout.setObjectName(u"numTimeout")
        self.numTimeout.setDecimals(1)
        self.numTimeout.setMinimum(1.000000000000000)
        self.numTimeout.setMaximum(5.000000000000000)
        self.numTimeout.setSingleStep(0.100000000000000)
        self.numTimeout.setValue(1.500000000000000)

        self.gridLayout_4.addWidget(self.numTimeout, 5, 1, 1, 1)

        self.label_5 = QLabel(self.isScan)
        self.label_5.setObjectName(u"label_5")

        self.gridLayout_4.addWidget(self.label_5, 4, 1, 1, 1)

        self.strUnPort = QLineEdit(self.isScan)
        self.strUnPort.setObjectName(u"strUnPort")

        self.gridLayout_4.addWidget(self.strUnPort, 3, 0, 1, 2)

        self.numThread = QSpinBox(self.isScan)
        self.numThread.setObjectName(u"numThread")
        self.numThread.setMaximum(6000)
        self.numThread.setSingleStep(100)
        self.numThread.setValue(400)

        self.gridLayout_4.addWidget(self.numThread, 5, 0, 1, 1)

        self.label_4 = QLabel(self.isScan)
        self.label_4.setObjectName(u"label_4")

        self.gridLayout_4.addWidget(self.label_4, 0, 0, 1, 1)

        self.label_6 = QLabel(self.isScan)
        self.label_6.setObjectName(u"label_6")

        self.gridLayout_4.addWidget(self.label_6, 2, 0, 1, 1)


        self.horizontalLayout_2.addWidget(self.isScan)

        self.isPair = QGroupBox(main)
        self.isPair.setObjectName(u"isPair")
        self.isPair.setMinimumSize(QSize(188, 0))
        self.isPair.setMaximumSize(QSize(192, 16777215))
        self.isPair.setCheckable(True)
        self.isPair.setChecked(False)
        self.verticalLayout_2 = QVBoxLayout(self.isPair)
        self.verticalLayout_2.setObjectName(u"verticalLayout_2")
        self.label_7 = QLabel(self.isPair)
        self.label_7.setObjectName(u"label_7")

        self.verticalLayout_2.addWidget(self.label_7)

        self.strPairPort = QLineEdit(self.isPair)
        self.strPairPort.setObjectName(u"strPairPort")

        self.verticalLayout_2.addWidget(self.strPairPort)

        self.label_8 = QLabel(self.isPair)
        self.label_8.setObjectName(u"label_8")

        self.verticalLayout_2.addWidget(self.label_8)

        self.strPairCode = QLineEdit(self.isPair)
        self.strPairCode.setObjectName(u"strPairCode")

        self.verticalLayout_2.addWidget(self.strPairCode)


        self.horizontalLayout_2.addWidget(self.isPair)


        self.gridLayout_2.addLayout(self.horizontalLayout_2, 2, 0, 1, 1)

        self.groupBox_4 = QGroupBox(main)
        self.groupBox_4.setObjectName(u"groupBox_4")
        self.gridLayout = QGridLayout(self.groupBox_4)
        self.gridLayout.setObjectName(u"gridLayout")
        self.isTop = QCheckBox(self.groupBox_4)
        self.isTop.setObjectName(u"isTop")

        self.gridLayout.addWidget(self.isTop, 0, 0, 1, 1)

        self.isSw = QCheckBox(self.groupBox_4)
        self.isSw.setObjectName(u"isSw")

        self.gridLayout.addWidget(self.isSw, 0, 1, 1, 1)

        self.isPoweroff = QCheckBox(self.groupBox_4)
        self.isPoweroff.setObjectName(u"isPoweroff")

        self.gridLayout.addWidget(self.isPoweroff, 0, 2, 1, 1)

        self.isRecord = QCheckBox(self.groupBox_4)
        self.isRecord.setObjectName(u"isRecord")
        self.isRecord.setStyleSheet(u"QCheckBox:checked{color: red;}")

        self.gridLayout.addWidget(self.isRecord, 0, 3, 1, 1)

        self.isHEVC = QCheckBox(self.groupBox_4)
        self.isHEVC.setObjectName(u"isHEVC")

        self.gridLayout.addWidget(self.isHEVC, 1, 0, 1, 1)

        self.isNoAudio = QCheckBox(self.groupBox_4)
        self.isNoAudio.setObjectName(u"isNoAudio")

        self.gridLayout.addWidget(self.isNoAudio, 1, 1, 1, 1)

        self.isNoMousehover = QCheckBox(self.groupBox_4)
        self.isNoMousehover.setObjectName(u"isNoMousehover")

        self.gridLayout.addWidget(self.isNoMousehover, 1, 2, 1, 1)

        self.isRecordAudio = QCheckBox(self.groupBox_4)
        self.isRecordAudio.setObjectName(u"isRecordAudio")

        self.gridLayout.addWidget(self.isRecordAudio, 1, 3, 1, 1)

        self.label_9 = QLabel(self.groupBox_4)
        self.label_9.setObjectName(u"label_9")

        self.gridLayout.addWidget(self.label_9, 2, 0, 1, 1)

        self.label_10 = QLabel(self.groupBox_4)
        self.label_10.setObjectName(u"label_10")

        self.gridLayout.addWidget(self.label_10, 2, 1, 1, 1)

        self.label_11 = QLabel(self.groupBox_4)
        self.label_11.setObjectName(u"label_11")

        self.gridLayout.addWidget(self.label_11, 2, 2, 1, 1)

        self.label_12 = QLabel(self.groupBox_4)
        self.label_12.setObjectName(u"label_12")

        self.gridLayout.addWidget(self.label_12, 2, 3, 1, 1)

        self.strKey = QLineEdit(self.groupBox_4)
        self.strKey.setObjectName(u"strKey")
        self.strKey.setMaximumSize(QSize(125, 16777215))

        self.gridLayout.addWidget(self.strKey, 3, 0, 1, 1)

        self.strSize = QLineEdit(self.groupBox_4)
        self.strSize.setObjectName(u"strSize")
        self.strSize.setMaximumSize(QSize(125, 16777215))

        self.gridLayout.addWidget(self.strSize, 3, 1, 1, 1)

        self.strBitrate = QLineEdit(self.groupBox_4)
        self.strBitrate.setObjectName(u"strBitrate")
        self.strBitrate.setMaximumSize(QSize(125, 16777215))

        self.gridLayout.addWidget(self.strBitrate, 3, 2, 1, 1)

        self.strBuffer = QLineEdit(self.groupBox_4)
        self.strBuffer.setObjectName(u"strBuffer")
        self.strBuffer.setMaximumSize(QSize(125, 16777215))

        self.gridLayout.addWidget(self.strBuffer, 3, 3, 1, 1)


        self.gridLayout_2.addWidget(self.groupBox_4, 3, 0, 1, 1)

        self.horizontalLayout_3 = QHBoxLayout()
        self.horizontalLayout_3.setObjectName(u"horizontalLayout_3")
        self.label_14 = QLabel(main)
        self.label_14.setObjectName(u"label_14")
        self.label_14.setOpenExternalLinks(True)

        self.horizontalLayout_3.addWidget(self.label_14)

        self.btnStart = QPushButton(main)
        self.btnStart.setObjectName(u"btnStart")
        sizePolicy = QSizePolicy(QSizePolicy.Policy.Minimum, QSizePolicy.Policy.Fixed)
        sizePolicy.setHorizontalStretch(0)
        sizePolicy.setVerticalStretch(0)
        sizePolicy.setHeightForWidth(self.btnStart.sizePolicy().hasHeightForWidth())
        self.btnStart.setSizePolicy(sizePolicy)
        font1 = QFont()
        font1.setFamilies([u"Microsoft YaHei UI"])
        font1.setPointSize(10)
        self.btnStart.setFont(font1)

        self.horizontalLayout_3.addWidget(self.btnStart)

        self.label_15 = QLabel(main)
        self.label_15.setObjectName(u"label_15")
        self.label_15.setAlignment(Qt.AlignmentFlag.AlignRight|Qt.AlignmentFlag.AlignTrailing|Qt.AlignmentFlag.AlignVCenter)

        self.horizontalLayout_3.addWidget(self.label_15)


        self.gridLayout_2.addLayout(self.horizontalLayout_3, 4, 0, 1, 1)

        QWidget.setTabOrder(self.btnStart, self.strPath)
        QWidget.setTabOrder(self.strPath, self.strIP)
        QWidget.setTabOrder(self.strIP, self.strPort)
        QWidget.setTabOrder(self.strPort, self.strKey)
        QWidget.setTabOrder(self.strKey, self.strSize)
        QWidget.setTabOrder(self.strSize, self.strBitrate)
        QWidget.setTabOrder(self.strBitrate, self.strBuffer)
        QWidget.setTabOrder(self.strBuffer, self.isTop)
        QWidget.setTabOrder(self.isTop, self.isHEVC)
        QWidget.setTabOrder(self.isHEVC, self.isSw)
        QWidget.setTabOrder(self.isSw, self.isNoAudio)
        QWidget.setTabOrder(self.isNoAudio, self.isPoweroff)
        QWidget.setTabOrder(self.isPoweroff, self.isRecord)
        QWidget.setTabOrder(self.isRecord, self.isScan)
        QWidget.setTabOrder(self.isScan, self.strPortRange)
        QWidget.setTabOrder(self.strPortRange, self.strUnPort)
        QWidget.setTabOrder(self.strUnPort, self.numThread)
        QWidget.setTabOrder(self.numThread, self.numTimeout)
        QWidget.setTabOrder(self.numTimeout, self.isPair)
        QWidget.setTabOrder(self.isPair, self.strPairPort)
        QWidget.setTabOrder(self.strPairPort, self.strPairCode)
        QWidget.setTabOrder(self.strPairCode, self.btnPath)

        self.retranslateUi(main)

        QMetaObject.connectSlotsByName(main)
    # setupUi

    def retranslateUi(self, main):
        main.setWindowTitle(QCoreApplication.translate("main", u"scrcpy GUI", None))
        self.label.setText(QCoreApplication.translate("main", u"scrcpy \u8def\u5f84:", None))
        self.strPath.setPlaceholderText(QCoreApplication.translate("main", u"\u53ef\u62d6\u653e\u76ee\u5f55\u6216\u6587\u4ef6, \u4e0e scrcpy \u540c\u76ee\u5f55\u53ef\u7559\u7a7a", None))
        self.btnPath.setText(QCoreApplication.translate("main", u"\u9009\u62e9\u76ee\u5f55", None))
        self.groupBox.setTitle(QCoreApplication.translate("main", u"\u8fde\u63a5", None))
        self.label_2.setText(QCoreApplication.translate("main", u"IP:", None))
        self.strIP.setInputMask(QCoreApplication.translate("main", u"000.000.000.000", None))
        self.label_3.setText(QCoreApplication.translate("main", u"\u7aef\u53e3:", None))
        self.strPort.setInputMask(QCoreApplication.translate("main", u"00000", None))
#if QT_CONFIG(tooltip)
        self.isScan.setToolTip(QCoreApplication.translate("main", u"\u542f\u52a8\u65f6\u626b\u63cf\u4e0a\u6b21\u8fde\u63a5ip\u7684\u6307\u5b9a\u8303\u56f4\u5185\u7aef\u53e3\n"
"(\u56e0\u65e0\u7ebf\u8c03\u8bd5\u6bcf\u6b21\u4f1a\u968f\u673a\u5f00\u7aef\u53e3)", None))
#endif // QT_CONFIG(tooltip)
        self.isScan.setTitle(QCoreApplication.translate("main", u"\u626b\u63cf", None))
#if QT_CONFIG(tooltip)
        self.label_13.setToolTip(QCoreApplication.translate("main", u"\u7ebf\u7a0b\u8d8a\u591a\u626b\u63cf\u8d8a\u5feb,\u4f46\u6570\u91cf\u4e0e\u5185\u5b58\u6709\u5173\n"
" 8G\u5185\u5b58: 800\n"
"16G\u5185\u5b58:1600", None))
#endif // QT_CONFIG(tooltip)
        self.label_13.setText(QCoreApplication.translate("main", u"\u7ebf\u7a0b:", None))
#if QT_CONFIG(tooltip)
        self.strPortRange.setToolTip(QCoreApplication.translate("main", u"\u4f8b: 36000-44000", None))
#endif // QT_CONFIG(tooltip)
        self.strPortRange.setInputMask(QCoreApplication.translate("main", u"00000-00000", None))
#if QT_CONFIG(tooltip)
        self.numTimeout.setToolTip(QCoreApplication.translate("main", u"\u8d8a\u5c0f\u901f\u5ea6\u8d8a\u5feb,\u4f46\u6f0f\u8fc7\u7684\u7aef\u53e3\u8d8a\u591a,\u4e0ewifi\u8fde\u63a5\u901f\u5ea6\u6709\u5173", None))
#endif // QT_CONFIG(tooltip)
#if QT_CONFIG(tooltip)
        self.label_5.setToolTip(QCoreApplication.translate("main", u"\u8d8a\u5c0f\u901f\u5ea6\u8d8a\u5feb,\u4f46\u6f0f\u8fc7\u7684\u7aef\u53e3\u8d8a\u591a,\u4e0ewifi\u8fde\u63a5\u901f\u5ea6\u6709\u5173", None))
#endif // QT_CONFIG(tooltip)
        self.label_5.setText(QCoreApplication.translate("main", u"\u8d85\u65f6:", None))
#if QT_CONFIG(tooltip)
        self.strUnPort.setToolTip(QCoreApplication.translate("main", u"\u626b\u63cf\u8303\u56f4\u5185\u975e\u8c03\u8bd5\u7528\u5e38\u5f00\u7aef\u53e3\n"
"\u591a\u4e2a\u7528\u7a7a\u683c\u5206\u9694", None))
#endif // QT_CONFIG(tooltip)
#if QT_CONFIG(tooltip)
        self.numThread.setToolTip(QCoreApplication.translate("main", u"\u7ebf\u7a0b\u8d8a\u591a\u626b\u63cf\u8d8a\u5feb,\u4f46\u6570\u91cf\u4e0e\u5185\u5b58\u6709\u5173\n"
" 8G\u5185\u5b58: 800\n"
"16G\u5185\u5b58:1600", None))
#endif // QT_CONFIG(tooltip)
#if QT_CONFIG(tooltip)
        self.label_4.setToolTip(QCoreApplication.translate("main", u"\u4f8b: 36000-44000", None))
#endif // QT_CONFIG(tooltip)
        self.label_4.setText(QCoreApplication.translate("main", u"\u7aef\u53e3\u8303\u56f4:", None))
#if QT_CONFIG(tooltip)
        self.label_6.setToolTip(QCoreApplication.translate("main", u"\u626b\u63cf\u8303\u56f4\u5185\u975e\u8c03\u8bd5\u7528\u5e38\u5f00\u7aef\u53e3\n"
"\u591a\u4e2a\u7528\u7a7a\u683c\u5206\u9694", None))
#endif // QT_CONFIG(tooltip)
        self.label_6.setText(QCoreApplication.translate("main", u"\u5c4f\u853d\u7aef\u53e3:", None))
#if QT_CONFIG(tooltip)
        self.isPair.setToolTip(QCoreApplication.translate("main", u"\u65e0\u7ebf\u8c03\u8bd5\u9996\u6b21\u8fde\u63a5\u9700\u914d\u5bf9\n"
"\u5177\u4f53\u770b\u624b\u673a(Android 11\u4ee5\u4e0a)\n"
"\u914d\u5bf9\u65f6ip,\u7aef\u53e3,\u914d\u5bf9\u7aef\u53e3,\u914d\u5bf9\u78014\u9879\u90fd\u8981\u586b", None))
#endif // QT_CONFIG(tooltip)
        self.isPair.setTitle(QCoreApplication.translate("main", u"\u914d\u5bf9", None))
        self.label_7.setText(QCoreApplication.translate("main", u"\u914d\u5bf9\u7aef\u53e3:", None))
        self.strPairPort.setInputMask(QCoreApplication.translate("main", u"00000", None))
        self.label_8.setText(QCoreApplication.translate("main", u"\u914d\u5bf9\u7801:", None))
        self.strPairCode.setInputMask(QCoreApplication.translate("main", u"000000", None))
        self.groupBox_4.setTitle(QCoreApplication.translate("main", u"\u9009\u9879", None))
        self.isTop.setText(QCoreApplication.translate("main", u"\u7f6e\u9876", None))
        self.isSw.setText(QCoreApplication.translate("main", u"\u7184\u5c4f\u4e14\u7981\u6b62\u4f11\u7720", None))
        self.isPoweroff.setText(QCoreApplication.translate("main", u"\u9000\u51fa\u540e\u9501\u5c4f", None))
#if QT_CONFIG(tooltip)
        self.isRecord.setToolTip(QCoreApplication.translate("main", u"\u5f55\u50cf\u6587\u4ef6\u4ee5\u65e5\u671f\u4fdd\u5b58\u5230 rec \u76ee\u5f55", None))
#endif // QT_CONFIG(tooltip)
        self.isRecord.setText(QCoreApplication.translate("main", u"\u5f55\u5c4f", None))
#if QT_CONFIG(tooltip)
        self.isHEVC.setToolTip(QCoreApplication.translate("main", u"scrcpy 2.0 \u4ee5\u4e0a\u652f\u6301", None))
#endif // QT_CONFIG(tooltip)
        self.isHEVC.setText(QCoreApplication.translate("main", u"H.265", None))
#if QT_CONFIG(tooltip)
        self.isNoAudio.setToolTip(QCoreApplication.translate("main", u"\u9ed8\u8ba4\u4e3a\u542f\u7528\u97f3\u9891\u8f6c\u53d1\n"
"\u9700 scrcpy 2.0 \u4e0e Android 11", None))
#endif // QT_CONFIG(tooltip)
        self.isNoAudio.setText(QCoreApplication.translate("main", u"\u7981\u7528\u97f3\u9891\u8f6c\u53d1", None))
#if QT_CONFIG(tooltip)
        self.isNoMousehover.setToolTip("")
#endif // QT_CONFIG(tooltip)
        self.isNoMousehover.setText(QCoreApplication.translate("main", u"\u7981\u7528\u9f20\u6807\u60ac\u505c", None))
#if QT_CONFIG(tooltip)
        self.isRecordAudio.setToolTip(QCoreApplication.translate("main", u"\u5f55\u97f3\u6587\u4ef6\u4ee5 AAC \u683c\u5f0f\u4fdd\u5b58\u5230 rec \u76ee\u5f55", None))
#endif // QT_CONFIG(tooltip)
        self.isRecordAudio.setText(QCoreApplication.translate("main", u"\u4ec5\u5f55\u97f3", None))
#if QT_CONFIG(tooltip)
        self.label_9.setToolTip(QCoreApplication.translate("main", u"\u9ed8\u8ba4\u4e3a\u5de6Alt\u952e: lalt\n"
"\u53ef\u9009\u503c: lctrl\u3001rctrl\u3001lalt\u3001ralt\u3001lsuper\u3001rsuper", None))
#endif // QT_CONFIG(tooltip)
        self.label_9.setText(QCoreApplication.translate("main", u"\u70ed\u952e:", None))
        self.label_10.setText(QCoreApplication.translate("main", u"\u9ad8\u5ea6:", None))
#if QT_CONFIG(tooltip)
        self.label_11.setToolTip(QCoreApplication.translate("main", u"\u9ed8\u8ba4\u7801\u7387\u662f 8 Mbps", None))
#endif // QT_CONFIG(tooltip)
        self.label_11.setText(QCoreApplication.translate("main", u"\u7801\u7387(Mbps):", None))
        self.label_12.setText(QCoreApplication.translate("main", u"\u7f13\u51b2(\u6beb\u79d2):", None))
#if QT_CONFIG(tooltip)
        self.strKey.setToolTip(QCoreApplication.translate("main", u"\u9ed8\u8ba4\u4e3a\u5de6Alt\u952e: lalt\n"
"\u53ef\u9009\u503c: lctrl\u3001rctrl\u3001lalt\u3001ralt\u3001lsuper\u3001rsuper", None))
#endif // QT_CONFIG(tooltip)
        self.strKey.setInputMask(QCoreApplication.translate("main", u"aaaaaa", None))
        self.strSize.setInputMask(QCoreApplication.translate("main", u"0000", None))
#if QT_CONFIG(tooltip)
        self.strBitrate.setToolTip(QCoreApplication.translate("main", u"\u9ed8\u8ba4\u7801\u7387\u662f 8 Mbps", None))
#endif // QT_CONFIG(tooltip)
        self.strBitrate.setInputMask(QCoreApplication.translate("main", u"0", None))
        self.strBuffer.setInputMask(QCoreApplication.translate("main", u"000", None))
        self.label_14.setText(QCoreApplication.translate("main", u"<a href=\"https://github.com/Genymobile/scrcpy\">scrcpy v3.2</a>", None))
        self.btnStart.setText(QCoreApplication.translate("main", u"\u542f \u52a8", None))
        self.label_15.setText(QCoreApplication.translate("main", u"by Courser", None))
    # retranslateUi

