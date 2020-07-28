# -*- coding: utf-8 -*-
"""
Created on Mon Jul 27 15:50:59 2020

@author: changhaizhao
"""

# -*- coding: utf-8 -*-
"""
Created on Fri Jul 24 13:55:31 2020

@author: changhaizhao
"""

import matplotlib.pyplot as plt
import pandas as pd
import datetime
import numpy as np
import math

#xlsx_file = "D:\工作目录\FROM潘博\LF_table(6).xlsx"
xlsx_file = "LF_table.xlsx"
deviceID = pd.read_excel(xlsx_file, usecols=[1])
dataTime = pd.read_excel(xlsx_file, usecols=[2])
dataIDTime = pd.read_excel(xlsx_file, usecols=[3])

deviceIDName=deviceID["device-id"].value_counts()


deviceNo1=deviceID[deviceID['device-id']=='510822208200LF03'].index.tolist()
#deviceNo1=deviceID[deviceID['device-id']=='120000000000LF02'].index.tolist()
#deviceNo1=deviceID[deviceID['device-id']=='510822208200LF01'].index.tolist()
#deviceNo1=deviceID[deviceID['device-id']=='510822208200LF02'].index.tolist()
#deviceNo1=deviceID[deviceID['device-id']=='510822208200LF04'].index.tolist()
#deviceNo1=deviceID[deviceID['device-id']=='120000000000LF01'].index.tolist()
#deviceNo1=deviceID[deviceID['device-id']=='510726216001LF01'].index.tolist()


deviceNo1DataTime=dataTime['time'][deviceNo1]
n_samples=len(deviceNo1DataTime)
deviceNo1DataMM=dataIDTime['mm'][deviceNo1]

OutTime=deviceNo1DataTime[1:n_samples]

index=0;
a=datetime.datetime(2019,12,31,23,59,59)
for b in OutTime:
    if a.__gt__(b):
        index += 1
    else:
        break

OutData=deviceNo1DataMM[1:n_samples]
temp=deviceNo1DataMM.iloc[index-1]
OutData= list(map(lambda x: x - temp, OutData)) 

OutTime=OutTime[index:-1]
OutData=OutData[index:-1]

nout=len(OutTime)
ProOutData=list([OutData[0]])
ProOutData.append(OutData[1])
ProOutData.append(OutData[2])
ProOutTime=list([OutTime.iloc[0]])
ProOutTime.append(OutTime.iloc[1])
ProOutTime.append(OutTime.iloc[2])
for i in range(3,nout):
    if ((OutData[i-1]+OutData[i-2]+OutData[i-3])/3) > (OutData[i]*1.5) :
        continue;
    ProOutData.append(OutData[i])
    ProOutTime.append(OutTime.iloc[i])

nPro=len(ProOutData)
for j  in range(0,nPro):
    if ProOutData[j]<0:
        ProOutData[j]=0.0

nout2=len(ProOutTime)
ProOutData2=list([ProOutData[0]])
ProOutData2.append(ProOutData[1])
ProOutData2.append(ProOutData[2])
ProOutTime2=list([ProOutTime[0]])
ProOutTime2.append(ProOutTime[1])
ProOutTime2.append(ProOutTime[2])
for i in range(3,nout2):
    if ((ProOutData[i-1]+ProOutData[i-2]+ProOutData[i-3])/3) > (ProOutData[i]*1.5) :
        continue;
    ProOutData2.append(ProOutData[i])
    ProOutTime2.append(ProOutTime[i])

plt.plot(ProOutTime2,ProOutData2,'b',label='111   Cumulative Displacement Curve')
plt.legend()
plt.xlabel("Date")
plt.ylabel("Cumulative Displacement/(mm)")
plt.show()

x = pd.DataFrame(ProOutData2)
x.to_excel('exam.xls')


def moving_average(interval, window_size):
    window = np.ones(int(window_size)) / float(window_size)
    return np.convolve(interval, window, 'same')  # numpy的卷积函数


y_av = moving_average(interval = ProOutData2, window_size = 10)
y_av2 = moving_average(interval = y_av, window_size = 10)
y_av3 = moving_average(interval = y_av2, window_size = 10)
plt.plot(ProOutTime2, y_av3, 'r',label='222  Smoothed Cumulative Displacement Curve')
plt.xlabel('Date')
plt.ylabel('Cumulative Displacement Curve')
plt.show()


y = pd.DataFrame(y_av3)
y.to_excel('examsmooth.xls')

#可以显示中文
plt.rcParams["font.sans-serif"] = ["SimHei"]
plt.rcParams['axes.unicode_minus'] = False

xxtemp = ProOutTime2[0]
xx = list(map(lambda x: (((x - xxtemp).seconds)/60/60/24 + (x - xxtemp).days), ProOutTime2)) 

P = np.polyfit(xx, y_av3, 16)    # 多项式函数拟合
poly_y = np.polyval(P, xx)

npoly=len(poly_y)
for i in range(1,npoly):
    if (poly_y[i]-poly_y[i-1])<0:
        ls = poly_y[i-1]-poly_y[i]
        for j in range(i,npoly):
            poly_y[j] = poly_y[j]+1.5*ls
            
plt.plot(xx, poly_y, 'k',label='333  处理后位移曲线')
plt.xlabel('Date')
plt.ylabel('Cumulative Displacement Curve')
plt.show()
            
            
#            
P1 = np.polyfit(xx, poly_y,20)    # 多项式函数拟合
poly_y1 = np.polyval(P1, xx)

plt.plot(xx, poly_y1, color="r", label="位移曲线")
plt.legend()
plt.show()


Pd1 = np.polyder(P1)  # 求导数
Pd1_y=np.polyval(Pd1,xx)
plt.plot(xx, Pd1_y, color="b", label="速度曲线")
plt.legend()
plt.show()

lsy=[]
v=1
lsn=len(Pd1_y)
for i in range(0,lsn):
    lsy.append(math.atan(Pd1_y[i])*180/v/3.14)


Pd2 = np.polyder(Pd1)  # 求导数
Pd2_y=np.polyval(Pd2,xx)
plt.plot(xx, Pd2_y, color="k", label="加速度曲线")
plt.legend()
plt.show()
