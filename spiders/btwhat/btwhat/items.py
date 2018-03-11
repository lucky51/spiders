# -*- coding: utf-8 -*-

# Define here the models for your scraped items
#
# See documentation in:
# https://doc.scrapy.org/en/latest/topics/items.html

import scrapy


class BtwhatItem(scrapy.Item):
    # define the fields for your item here like:
    # name = scrapy.Field()
    file_type = scrapy.Field()
    file_createtime = scrapy.Field()
    file_hot = scrapy.Field()
    file_size = scrapy.Field()
    file_count = scrapy.Field()
    file_link = scrapy.Field()
    file_name = scrapy.Field()
    file_url = scrapy.Field()