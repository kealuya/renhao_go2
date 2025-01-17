package excel

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nguyenthenguyen/docx"
	"github.com/xuri/excelize/v2"
)

// Employee 结构体用于存储员工信息
type Employee struct {
	ID                 string  // 工号
	Name               string  // 姓名
	Overtime           float64 // 加班天数
	WorkScore          float64 // 工作成果评分
	SkillScore         float64 // 专业技能提升评分
	AttitudeScore      float64 // 工作态度评分
	LeaveCount         float64 // 请假天数
	ComprehensiveScore float64 // 综合评分
	Evaluation         string  // 评价
}

func Go_excel() {
	// tongjijiaban()
	// tongjiQingjia()
	err := GeneratePerformanceReviews()
	if err != nil {
		log.Fatalf("生成绩效考评确认书时出错: %v", err)
	}
}

func tongjijiaban() {
	// 存储汇总数据的map
	summaryMap := make(map[string]*Employee)

	// 指定要处理的文件夹路径
	folderPath := "/Users/kealuya/办公/办公文稿/2024浩天教育考勤/自动化/jiaban/" // 根据实际情况修改路径

	// 遍历文件夹下的所有Excel文件
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理xlsx文件
		if !info.IsDir() && (filepath.Ext(path) == ".xlsx" || filepath.Ext(path) == ".xls") {
			err := processExcelFile(path, summaryMap)
			if err != nil {
				log.Printf("处理文件 %s 时出错: %v\n", path, err)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("遍历文件夹出错: %v", err)
	}

	// 创建新的Excel文件保存汇总结果
	err = createSummaryExcel(summaryMap)
	if err != nil {
		log.Fatalf("创建汇总文件时出错: %v", err)
	}

	fmt.Println("数据汇总完成！")
}

// 处理单个Excel文件
func processExcelFile(filePath string, summaryMap map[string]*Employee) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer f.Close()

	// 获取第一个sheet的名称
	firstSheet := f.GetSheetName(0)
	if firstSheet == "" {
		return fmt.Errorf("未找到任何sheet")
	}

	// 读取第一个sheet的所有行
	rows, err := f.GetRows(firstSheet)
	if err != nil {
		return fmt.Errorf("读取Sheet1失败: %v", err)
	}

	// 从第二行开始处理数据（假设第一行是标题）
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 3 {
			continue // 跳过数据不完整的行
		}

		id := row[0]
		name := row[1]
		overtime := 0.0
		fmt.Sscanf(row[2], "%f", &overtime)

		// 更新或创建员工记录
		if emp, exists := summaryMap[id]; exists {
			emp.Overtime += overtime
		} else {
			summaryMap[id] = &Employee{
				ID:       id,
				Name:     name,
				Overtime: overtime,
			}
		}
	}

	return nil
}

// 创建汇总Excel文件
func createSummaryExcel(summaryMap map[string]*Employee) error {
	f := excelize.NewFile()
	defer f.Close()

	// 设置标题行
	f.SetCellValue("Sheet1", "A1", "工号")
	f.SetCellValue("Sheet1", "B1", "姓名")
	f.SetCellValue("Sheet1", "C1", "加班天数")

	// 写入数据
	row := 2
	for _, emp := range summaryMap {
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), emp.ID)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), emp.Name)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), emp.Overtime)
		row++
	}

	// 保存文件
	if err := f.SaveAs("汇总结果.xlsx"); err != nil {
		return fmt.Errorf("保存汇总文件失败: %v", err)
	}

	return nil
}

// 处理请假情况的汇总
func tongjiQingjia() {
	// 存储汇总数据的map，key是工号，value是一个包含12个月请假情况的map
	summaryMap := make(map[string]map[int]string)

	// 指定要处理的文件夹路径
	folderPath := "/Users/kealuya/办公/办公文稿/2024浩天教育考勤/自动化/qingjia/" // 根据实际情况修改路径

	// 处理每个月份的文件
	for month := 1; month <= 12; month++ {
		fileName := fmt.Sprintf("考勤工作表2024年%d月_公式.xlsx", month)
		filePath := filepath.Join(folderPath, fileName)

		// 检查文件是否存在
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			continue // 如果文件不存在，跳过这个月份
		}

		// 处理文件
		err := processQingjiaFile(filePath, month, summaryMap)
		if err != nil {
			log.Printf("处理%d月请假文件时出错: %v\n", month, err)
		}
	}

	// 创建汇总文件
	err := createQingjiaSummaryExcel(summaryMap)
	if err != nil {
		log.Fatalf("创建请假汇总文件时出错: %v", err)
	}

	fmt.Println("请假数据汇总完成！")
}

// 处理单个请假Excel文件
func processQingjiaFile(filePath string, month int, summaryMap map[string]map[int]string) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer f.Close()

	// 获取Sheet1
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return fmt.Errorf("读取Sheet1失败: %v", err)
	}

	// 从第二行开始处理数据（假设第一行是标题）
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 3 {
			continue // 跳过数据不完整的行
		}

		id := row[0]      // A列工号
		qingjia := row[2] // C列请假情况

		// 如果这个工号还没有对应的map，创建一个新的
		if _, exists := summaryMap[id]; !exists {
			summaryMap[id] = make(map[int]string)
		}

		// 保存这个月的请假情况
		summaryMap[id][month] = qingjia
	}

	return nil
}

// 创建请假汇总Excel文件
func createQingjiaSummaryExcel(summaryMap map[string]map[int]string) error {
	f := excelize.NewFile()
	defer f.Close()

	// 设置标题行
	f.SetCellValue("Sheet1", "A1", "工号")
	for month := 1; month <= 12; month++ {
		col := string(rune('B' + month - 1))
		f.SetCellValue("Sheet1", fmt.Sprintf("%s1", col), fmt.Sprintf("%d月份", month))
	}

	// 写入数据
	row := 2
	for id, monthData := range summaryMap {
		// 写入工号
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), id)

		// 写入每个月的请假情况
		for month := 1; month <= 12; month++ {
			col := string(rune('B' + month - 1))
			qingjia := monthData[month] // 如果没有数据，会返回空字符串
			f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", col, row), qingjia)
		}
		row++
	}

	// 保存文件
	if err := f.SaveAs("/Users/kealuya/办公/办公文稿/2024浩天教育考勤/自动化/qingjia/请假汇总结果.xlsx"); err != nil {
		return fmt.Errorf("保存请假汇总文件失败: %v", err)
	}

	return nil
}

// GeneratePerformanceReviews 生成绩效考评确认书
func GeneratePerformanceReviews() error {
	// 打开Excel文件
	f, err := excelize.OpenFile("/Users/kealuya/办公/办公文稿/2024浩天教育考勤/自动化/绩效确认书/浩天教育2024_人员评价表.xlsx")
	if err != nil {
		return fmt.Errorf("打开Excel文件失败: %v", err)
	}
	defer f.Close()

	// 获取"人员评价"sheet的所有行
	rows, err := f.GetRows("人员评价")
	if err != nil {
		return fmt.Errorf("读取人员评价sheet失败: %v", err)
	}

	// 从第3行开始处理数据
	for i := 2; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 10 || row[2] == "" { // 检查是否有足够的列数和姓名是否为空
			continue
		}

		emp := Employee{
			Name:               row[2],             // C列：姓名
			WorkScore:          parseFloat(row[3]), // D列：工作成果评分
			SkillScore:         parseFloat(row[4]), // E列：专业技能提升评分
			AttitudeScore:      parseFloat(row[5]), // F列：工作态度评分
			Overtime:           parseFloat(row[6]), // G列：加班天数
			LeaveCount:         parseFloat(row[7]), // H列：请假天数
			ComprehensiveScore: parseFloat(row[8]), // I列：综合评分
			Evaluation:         row[9],             // J列：评价
		}

		// 为每个员工生成绩效考评确认书
		err := generateWordDocument(emp)
		if err != nil {
			log.Printf("生成%s的绩效考评确认书时出错: %v\n", emp.Name, err)
		}
	}

	return nil
}

// parseFloat 安全地将字符串转换为float64
func parseFloat(s string) float64 {
	var result float64
	fmt.Sscanf(s, "%f", &result)
	return result
}

// generateWordDocument 根据模板生成Word文档
func generateWordDocument(emp Employee) error {
	templatePath := "/Users/kealuya/办公/办公文稿/2024浩天教育考勤/自动化/绩效确认书/绩效考评确认书-模板.docx"
	outputPath := fmt.Sprintf("/Users/kealuya/办公/办公文稿/2024浩天教育考勤/自动化/绩效确认书/output/%s-绩效考评确认书.docx", emp.Name)

	// 读取模板文件
	r, err := docx.ReadDocxFile(templatePath)
	if err != nil {
		return fmt.Errorf("读取Word模板失败: %v", err)
	}
	defer r.Close()

	// 获取文档
	doc := r.Editable()

	// 准备替换的数据并进行替换
	doc.Replace("***Name***", emp.Name, -1)
	doc.Replace("***WorkScore***", fmt.Sprintf("%.1f", emp.WorkScore), -1)
	doc.Replace("***SkillScore***", fmt.Sprintf("%.1f", emp.SkillScore), -1)
	doc.Replace("***AttitudeScore***", fmt.Sprintf("%.1f", emp.AttitudeScore), -1)
	doc.Replace("***Overtime***", fmt.Sprintf("%.1f", emp.Overtime), -1)
	doc.Replace("***LeaveCount***", fmt.Sprintf("%.1f", emp.LeaveCount), -1)
	doc.Replace("***ComprehensiveScore***", fmt.Sprintf("%.1f", emp.ComprehensiveScore), -1)
	doc.Replace("***Evaluation***", emp.Evaluation, -1)

	// 保存新文件
	err = doc.WriteToFile(outputPath)
	if err != nil {
		return fmt.Errorf("保存Word文件失败: %v", err)
	}

	fmt.Printf("已生成%s的绩效考评确认书\n", emp.Name)
	return nil
}
