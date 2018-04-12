#!/usr/bin/python
# -*- coding: UTF-8 -*-
# _author:lucky51
# date:2018/4/11
import requests
from lxml import etree
import time
import os
import sys
import collections
reload(sys)
sys.setdefaultencoding("utf-8")


class SpOverItem(object):
    """爬取qhrb排行榜数据,只实践了默认"""
    item_run = True
    base_path = os.path.dirname(os.path.abspath(__file__))
    headers = {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36',
        'Host': 'spds.qhrb.com.cn',
        'Content-Type': 'application/x-www-form-urlencoded'
    }

    def __init__(self, starturl):
        self.url = starturl
        self.offset = 1
        self.request_data = requests.session()

    def __get_form_str(self, form=None):
        if form is not None:
            return '&'.join([k+'='+v for k, v in form.items()])
        else:
            return ''

    @staticmethod
    def __save_local_file(path, text):

        # text += "\r\n"
        # text += "=" * 200
        fil_name = str(int(time.time())) + '.txt'
        file_name = os.path.join(path, fil_name)
        with open(file_name, 'w') as f:
            f.write(text)

    @staticmethod
    def __check_lis_empty(lis):
        if lis is None:
            return ""
        else:
            if len(lis) > 0:
                return lis[0]
        return ""

    def __get_response(self, first=True, data=None, func=None):
        if data is not None:
            print data
        if first:
            self.response = self.request_data.get(self.url, headers=self.headers)
            self.page_count = int(etree.HTML(self.response.text).xpath('//*[@id="AspNetPager1_input"]/option[last()]/@value')[0])
        else:
            self.prev_reponse = self.response
            self.response = self.request_data.post(self.url, headers=self.headers, data=data)
        self.response_selector = etree.HTML(self.response.text)

        self.totals = self.response_selector.xpath('//*[@id="AspNetPager1_input"]/option[last()]/@value')
        if len(self.totals) > 0:
            self.prev_reponse = self.response
            self.offset += 1
            self.res_trs = self.response_selector.xpath('//*[@id="form1"]/div[6]/table//tr')
            save_text = ''
            for i in self.res_trs:
                save_text += (''.join(i.xpath('./td/text()')).replace(r'\n', ' ').replace(r'\r', ' '))
            self.__save_local_file(self.base_path, save_text)
        else:
            self.response = self.prev_reponse
        print str(self.offset) + "页码"
        if self.page_count == self.offset:
            return False
        else:
            return True

    def request(self):
        while self.item_run:
            time.sleep(1)
            if self.offset == 1:
                self.item_run = self.__get_response()
            else:
                self.__get_hidden_input(offset=self.offset)
                self.item_run = self.__get_response(data=self.__get_hidden_input(self.offset), first=False)

    def __get_hidden_input(self, offset=1): # txtTradeDate
        event_taget = "AspNetPager1"
        self.current_response_selecotr = etree.HTML(self.response.text)
        view_state = self.__check_lis_empty(self.current_response_selecotr.xpath('//*[@id="__VIEWSTATE"]/@value'))
        event_targument = "" # str(offset)
        view_state_generator = self.__check_lis_empty(
            self.current_response_selecotr.xpath('//*[@id="__VIEWSTATEGENERATOR"]/@value')

        )

        event_validation = self.__check_lis_empty(self.current_response_selecotr.xpath('//*[@id="__EVENTVALIDATION"]/@value'))
        # print event_validation
        hid_account_type = self.__check_lis_empty(self.current_response_selecotr.xpath('//*[@id="hidAccountType"]/@value'))
        hid_match_type = self.__check_lis_empty(self.current_response_selecotr.xpath('//*[@id="hidMatchType"]/@value'))
        # hid_trade_breed_type = self.__check_lis_empty(self.current_response_selecotr.xpath('//*[@id="hidTradeBreedType"]/@value'))
        hid_rank_type = self.__check_lis_empty(self.current_response_selecotr.xpath('//*[@id="hidRankType"]/@value'))
        # 资金账号
        txt_internal_account = ""
        # 昵称 self.response_selector.xpath('//*[@id="txtNickName"]/@value')[0]
        txt_nick_name = ""
        # 交易日期
        txt_trade_date = self.__check_lis_empty(self.current_response_selecotr.xpath('//*[@id="txtTradeDate"]/@value'))
        dic = collections.OrderedDict()
        dic['__EVENTTARGET'] = event_taget
        dic['__EVENTARGUMENT'] = event_targument
        dic['__VIEWSTATE'] = view_state
        dic['__VIEWSTATEGENERATOR'] = view_state_generator
        dic['__EVENTVALIDATION'] = event_validation
        dic['hidAccountType'] = hid_account_type
        dic['hidMatchType'] = hid_match_type
        dic['hidTradeBreedType'] = 0
        dic['hidRankType'] = hid_rank_type
        dic['txtInternalAccount'] = txt_internal_account
        dic['txtNickName'] = txt_nick_name
        dic['txtTradeDate'] = txt_trade_date
        dic['AspNetPager1_input'] = str(offset-1)
        return dic


if __name__ == '__main__':
    startUrl = 'http://spds.qhrb.com.cn/SP12/SPOverSee1.aspx'
    sp_item = SpOverItem(startUrl)
    sp_item.request()
