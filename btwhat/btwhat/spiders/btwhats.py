# -*- coding: utf-8 -*-
import scrapy
from scrapy.linkextractors import LinkExtractor
from scrapy.spiders import CrawlSpider, Rule
import urllib
from btwhat.items import BtwhatItem
import sys

reload(sys)

sys.setdefaultencoding('utf8')

class BtwhatsSpider(CrawlSpider):
    name = 'btwhats'
    allowed_domains = ['btwhat.info']
    def __init__(self, key_word='', *args, **kwargs):
        super(BtwhatsSpider, self).__init__(*args, **kwargs)
        self.key_words = key_word
        quote_str = urllib.quote(self.key_words)
        # 网址就不搞出来啦
        zero_url = 'http://www.btwhat.info/search/' + quote_str + '.html'
        self.start_urls = [zero_url]


    rules = (
        Rule(LinkExtractor(allow=r'\/search\/b-[\s\S]*\.html'),callback='root_url', follow=True),
         Rule(   LinkExtractor(
                allow=r'\/search\/b-[a-z,A-Z,0-9]+\/[0-9]+-[0-9]+\.html'), callback='content_url', follow=True
            ),
        Rule(LinkExtractor(allow=r'\/wiki\/.*\.html'), callback='parse_item', follow=False)
    )

    def root_url(self, response):
        pass


    def content_url(self, response):
        pass


    def parse_item(self, response):
        i = BtwhatItem()
        script_txt  = response.xpath('//*[@id="wall"]/h2/script/text()').extract()
        if len(script_txt) !=0:
            url_str = script_txt[0].replace('document.write(decodeURIComponent(', '').replace('));', '').replace('"','')
            link_name = urllib.unquote(str(url_str.replace('+', '')))
            i["file_name"] = link_name
            print "*" * 10
            #print link_name
            print "*" * 10
        file_nodes = response.xpath('//*[@id="wall"]/div/table/tr[last()]/td/text()').extract()
        print "#" * 10
        print file_nodes
        print "#" * 10
        if len(file_nodes) > 0 :
            i["file_type"] = file_nodes[0].replace('\n', '')
            i["file_createtime"] = file_nodes[1].replace('\n', '')
            i["file_hot"] = file_nodes[2].replace('\n', '')
            i["file_size"] = file_nodes[3].replace('\n', '')
        i["file_url"] = response.url
        file_link = response.xpath('//*[@id="wall"]/div[1]/div[1]/div[2]/a/@href').extract()
        if len(file_link) > 0:
            i["file_link"] = file_link[0]
        yield i
