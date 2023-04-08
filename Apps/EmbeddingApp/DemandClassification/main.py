import os
import sys
import openai
import configparser
from classification import Classification


def setOpenaiKey():   
    config_file = configparser.ConfigParser()     
    config_file.read("./././configurations.ini")        
    openai.api_key = config_file['OpenaiSettings']['api_key']

def main():
    setOpenaiKey()
    c = Classification()
    source_file = ["MX,SP微商城运维支持","敦奴导购中心的搭建","成品采购单的打印简易模板，把总数量下的交付日期去除"]
    c.find_print_related_requirements(source_file)

if __name__ == "__main__":
    main()  