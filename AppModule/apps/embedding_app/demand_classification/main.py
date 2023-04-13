from .classification import Classification

def main(question:str):    
    c = Classification()
    source_file = ["MX,SP微商城运维支持","敦奴导购中心的搭建","成品采购单的打印简易模板，把总数量下的交付日期去除"]
    c.find_print_related_requirements(source_file)
