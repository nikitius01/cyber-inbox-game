package tasks

import (
	"fmt"
	"time"
)

type seedCase struct {
	name       string
	email      string
	subject    string
	body       string
	phishing   bool
	link       Link
	attachment Attachment
	flags      []RedFlag
	explain    string
}

func BuildSeedTasks() []Task {
	tasks := []Task{}
	for _, category := range []Difficulty{Easy, Medium, Hard, Nightmare} {
		for i, item := range categoryCases(category) {
			id := fmt.Sprintf("%s_%02d", string(category), i+1)
			tasks = append(tasks, buildTask(id, category, item))
		}
	}
	return tasks
}

func buildTask(id string, category Difficulty, item seedCase) Task {
	headers := map[string]any{
		"Return-Path": fmt.Sprintf("<%s>", item.email),
		"From":        fmt.Sprintf("%s <%s>", item.name, item.email),
		"Reply-To":    item.email,
		"To":          "student@example.edu",
		"Subject":     item.subject,
		"Date":        "Tue, 8 Jul 2026 10:42:12 +0300",
		"Message-ID":  fmt.Sprintf("<%s@mail.training.local>", id),
		"Received": []string{
			"from mx.training.local by inbox.training.local",
			fmt.Sprintf("from smtp.%s by mx.training.local", domainOf(item.email)),
		},
		"SPF":   "pass",
		"DKIM":  "pass",
		"DMARC": "pass",
	}
	if item.phishing {
		headers["Reply-To"] = "verify@temporary-check.example"
		headers["SPF"] = "fail"
		headers["DKIM"] = "none"
		headers["DMARC"] = "fail"
	}

	return Task{
		ID:          id,
		Category:    category,
		SenderName:  item.name,
		SenderEmail: item.email,
		Subject:     item.subject,
		Body:        item.body,
		IsPhishing:  item.phishing,
		Links:       nonEmptyLinks(item.link),
		Attachments: nonEmptyAttachments(item.attachment),
		Raw: RawEmail{
			Headers:  headers,
			BodyText: item.body,
			BodyHTML: fmt.Sprintf("<p>%s</p>", item.body),
			Source:   fmt.Sprintf("From: %s <%s>\nTo: student@example.edu\nSubject: %s\nSPF: %v\nDKIM: %v\nDMARC: %v\n\n%s", item.name, item.email, item.subject, headers["SPF"], headers["DKIM"], headers["DMARC"], item.body),
		},
		RedFlags:    item.flags,
		Explanation: item.explain,
		CreatedAt:   time.Now().UTC(),
	}
}

func nonEmptyLinks(link Link) []Link {
	if link.ActualURL == "" {
		return []Link{}
	}
	return []Link{link}
}

func nonEmptyAttachments(attachment Attachment) []Attachment {
	if attachment.FileName == "" {
		return []Attachment{}
	}
	return []Attachment{attachment}
}

func domainOf(email string) string {
	for i, char := range email {
		if char == '@' {
			return email[i+1:]
		}
	}
	return "unknown.local"
}

func categoryCases(category Difficulty) []seedCase {
	brands := []string{"PayPal", "Google Workspace", "Microsoft 365", "Steam", "UniBank", "GitHub", "Dropbox", "Netflix", "Delivery Club", "University Portal"}
	cases := make([]seedCase, 0, 30)
	for i := 0; i < 30; i++ {
		brand := brands[i%len(brands)]
		phishing := i%3 != 1
		legitDomain := map[Difficulty]string{Easy: "example.edu", Medium: "notifications.example.edu", Hard: "corp.example.edu", Nightmare: "sso.example.edu"}[category]
		badDomain := map[Difficulty]string{Easy: "security-free-prize.example", Medium: "login-confirm.example", Hard: "examp1e-edu.example", Nightmare: "sso-example-edu.com"}[category]
		senderDomain := legitDomain
		if phishing {
			senderDomain = badDomain
		}
		subject := subjectFor(category, brand, i, phishing)
		body := legitimateBody(category, brand, i)
		flags := []RedFlag{}
		explain := "Письмо легитимное: домен отправителя ожидаемый, проверки SPF/DKIM/DMARC проходят, нет опасных вложений или скрытых ссылок."
		link := Link{
			VisibleText:  fmt.Sprintf("https://%s/account", legitDomain),
			ActualURL:    fmt.Sprintf("https://%s/account", legitDomain),
			Domain:       legitDomain,
			Protocol:     "https",
			IsSuspicious: false,
		}
		attachment := Attachment{}
		if phishing {
			body = phishingBody(category, brand, i)
			link = suspiciousLink(category, legitDomain, badDomain)
			attachment = suspiciousAttachment(category, i)
			flags = baseFlags(category, badDomain, attachment.FileName)
			explain = "Письмо опасное: адрес и технические заголовки не совпадают с ожидаемым доменом, ссылка ведет на сторонний ресурс, а вложение или тон письма создают дополнительный риск."
		}
		cases = append(cases, seedCase{
			name:       brand + " Support",
			email:      fmt.Sprintf("support@%s", senderDomain),
			subject:    subject,
			body:       body,
			phishing:   phishing,
			link:       link,
			attachment: attachment,
			flags:      flags,
			explain:    explain,
		})
	}
	return cases
}

func subjectFor(category Difficulty, brand string, index int, phishing bool) string {
	if phishing {
		switch category {
		case Easy:
			return fmt.Sprintf("%s: срочное подтверждение доступа #%02d", brand, index+1)
		case Medium:
			return fmt.Sprintf("%s: проверка необычной активности в аккаунте #%02d", brand, index+1)
		case Hard:
			return fmt.Sprintf("%s: обновление политики доступа для рабочей учетной записи #%02d", brand, index+1)
		default:
			return fmt.Sprintf("%s: эскалация инцидента безопасности INC-%04d", brand, 7000+index)
		}
	}
	switch category {
	case Easy:
		return fmt.Sprintf("%s: обычное уведомление по аккаунту #%02d", brand, index+1)
	case Medium:
		return fmt.Sprintf("%s: отчет о плановом изменении настроек #%02d", brand, index+1)
	case Hard:
		return fmt.Sprintf("%s: подтверждение завершенной заявки CHG-%04d", brand, 1200+index)
	default:
		return fmt.Sprintf("%s: сводка аудита доступа REF-%04d", brand, 4300+index)
	}
}

func legitimateBody(category Difficulty, brand string, index int) string {
	switch category {
	case Easy:
		return fmt.Sprintf(`Здравствуйте.

Это обычное информационное уведомление от %s. В вашем аккаунте завершена плановая проверка контактных данных. Никаких действий с вашей стороны не требуется.

Если вы не узнаете это действие, откройте приложение или официальный сайт вручную через закладку браузера. Мы не просим сообщать пароль, одноразовые коды или данные банковской карты в ответ на это письмо.

С уважением,
служба уведомлений %s
Номер события: SAFE-%04d`, brand, brand, 1000+index)
	case Medium:
		return fmt.Sprintf(`Здравствуйте.

Команда %s сообщает, что для вашей учетной записи была применена плановая настройка безопасности: обновлены параметры уведомлений и список доверенных устройств. Изменение выполнено в рамках регулярного обслуживания.

Для проверки деталей используйте официальный портал организации или приложение %s. Ссылки в этом письме ведут на основной домен сервиса, а вложений, требующих запуска макросов или установки программ, нет.

Если вы не ожидали такого уведомления, создайте обращение через внутренний helpdesk. Не пересылайте письмо третьим лицам, потому что оно содержит служебный номер операции.

Справочный номер: OK-%04d`, brand, brand, 2100+index)
	case Hard:
		return fmt.Sprintf(`Добрый день.

По заявке CHG-%04d завершена проверка настроек %s для вашей рабочей учетной записи. Изменение было согласовано администратором отдела и выполнено в стандартное окно обслуживания.

Письмо носит уведомительный характер. Для просмотра журнала действий откройте корпоративный портал напрямую, используя сохраненную ссылку или адрес из документации. Поддержка %s не запрашивает пароль, seed-фразы, одноразовые коды и не просит отключать защитные расширения браузера.

Техническая сводка:
- тип события: плановое обслуживание;
- статус: завершено;
- вложения: отсутствуют;
- дальнейшие действия: не требуются.

С уважением,
служба сопровождения %s`, 1200+index, brand, brand, brand)
	default:
		return fmt.Sprintf(`Коллеги, добрый день.

Отправляем сводку по аудиту доступа REF-%04d для сервиса %s. Проверка показала, что ваша учетная запись используется в рамках ожидаемого профиля: входы происходили из известных сетей, а последние изменения ролей совпадают с утвержденной заявкой.

Пожалуйста, не отвечайте на это письмо с персональными данными. Если нужен полный отчет, откройте раздел аудита во внутреннем портале вручную. В целях безопасности мы не прикладываем исполняемые файлы, архивы с паролем или HTML-формы авторизации.

Контрольные признаки легитимности:
- адрес отправителя находится в ожидаемом домене;
- Reply-To совпадает с отправителем;
- SPF, DKIM и DMARC проходят проверку;
- ссылка ведет на известный домен организации.

Отдел информационной безопасности`, 4300+index, brand)
	}
}

func phishingBody(category Difficulty, brand string, index int) string {
	switch category {
	case Easy:
		return fmt.Sprintf(`СРОЧНО!

Ваш аккаунт %s будет ограничен через 10 минут из-за неподтвержденной проверки безопасности. Чтобы избежать блокировки, перейдите по ссылке ниже и подтвердите пароль, номер телефона и резервный email.

Если форма не откроется, временно отключите блокировщик рекламы или попробуйте другой браузер. Не закрывайте страницу до завершения проверки, иначе заявка будет автоматически отклонена.

Код предупреждения: ALERT-%04d
Служба безопасности %s`, brand, 5100+index, brand)
	case Medium:
		return fmt.Sprintf(`Здравствуйте.

Мы обнаружили необычный вход в %s с нового устройства. Для предотвращения ограничений подтвердите учетную запись до конца рабочего дня. Проверка займет меньше минуты и необходима всем пользователям, получившим это уведомление.

Перейдите по ссылке "официального портала" и введите текущий пароль. После подтверждения система предложит обновить резервный способ входа. Если вы уже проходили проверку ранее, повторите ее еще раз, потому что предыдущая сессия могла быть не сохранена.

Обратите внимание: письмо отправлено автоматически, ответы на него не обрабатываются. При возникновении ошибки используйте вложенный файл или ссылку из письма.

Номер проверки: SEC-%04d`, brand, 6200+index)
	case Hard:
		return fmt.Sprintf(`Добрый день.

В рамках обновления политики доступа %s требуется повторная авторизация сотрудников и студентов, у которых включены расширенные права. Процедура выглядит как стандартная проверка SSO и занимает 2-3 минуты.

Пожалуйста, используйте ссылку ниже. Она может отличаться от привычного адреса портала, потому что трафик временно переведен на резервный шлюз. После входа подтвердите пароль и одноразовый код, чтобы система могла синхронизировать настройки. Если браузер покажет предупреждение о сертификате, выберите продолжение.

Детали изменения:
- окно применения: сегодня;
- причина: обновление матрицы ролей;
- влияние: возможная блокировка доступа к материалам курса;
- действие пользователя: обязательная повторная авторизация.

Команда сопровождения %s`, brand, brand)
	default:
		return fmt.Sprintf(`Коллеги, добрый день.

По инциденту INC-%04d в сервисе %s требуется срочная верификация владельца учетной записи. Мы зафиксировали цепочку входов, похожую на компрометацию, и временно переносим пользователей в отдельную группу проверки.

Для закрытия инцидента откройте приложенный отчет или перейдите по ссылке из письма. Форма запросит корпоративный логин, пароль, одноразовый код и подтверждение текущего устройства. Это необходимо, чтобы избежать отключения доступа к почте, учебному порталу и файловому хранилищу.

Обратите внимание на нестандартный порядок: проверка выполняется не через основной портал, а через резервный домен. Такой порядок введен временно из-за нагрузки на SSO. Не пересылайте письмо в общий чат, чтобы не задерживать обработку вашей учетной записи.

С уважением,
оперативная группа реагирования %s`, 7000+index, brand, brand)
	}
}

func suspiciousLink(category Difficulty, legitDomain, badDomain string) Link {
	protocol := "https"
	if category == Easy {
		protocol = "http"
	}
	return Link{
		VisibleText:  fmt.Sprintf("https://%s/account", legitDomain),
		ActualURL:    fmt.Sprintf("%s://%s/secure-login", protocol, badDomain),
		Domain:       badDomain,
		Protocol:     protocol,
		IsShortened:  category == Nightmare,
		IsSuspicious: true,
		RiskReason:   "Видимый текст ссылки отличается от реального домена назначения.",
	}
}

func suspiciousAttachment(category Difficulty, index int) Attachment {
	if index%2 == 0 {
		return Attachment{}
	}
	file := map[Difficulty]string{
		Easy:      "gift_card.pdf.exe",
		Medium:    "invoice_approved.xlsm",
		Hard:      "security_update.docm",
		Nightmare: "incident_report.zip",
	}[category]
	return Attachment{
		FileName:     file,
		DisplayName:  file,
		Extension:    file[stringsLastIndex(file, ".")+1:],
		MimeType:     "application/octet-stream",
		SizeKB:       128 + index*7,
		Hash:         fmt.Sprintf("simulated_sha256_%02d_%s", index, category),
		IsSuspicious: true,
		RiskReason:   "Тип файла часто используется для доставки вредоносных макросов или исполняемого кода.",
	}
}

func stringsLastIndex(value, needle string) int {
	last := 0
	for i := range value {
		if value[i:i+1] == needle {
			last = i
		}
	}
	return last
}

func baseFlags(category Difficulty, badDomain, attachment string) []RedFlag {
	flags := []RedFlag{
		{Type: "header", Field: "SPF", Value: "fail", Explanation: "Домен отправителя не прошел SPF-проверку."},
		{Type: "header", Field: "DMARC", Value: "fail", Explanation: "DMARC не подтверждает подлинность отправителя."},
		{Type: "link", Field: "actualUrl", Value: badDomain, Explanation: "Реальная ссылка ведет не на ожидаемый домен сервиса."},
	}
	if category == Easy {
		flags = append(flags, RedFlag{Type: "content", Field: "body", Value: "СРОЧНО", Explanation: "Давление временем часто используется в фишинге."})
	}
	if attachment != "" {
		flags = append(flags, RedFlag{Type: "attachment", Field: "fileName", Value: attachment, Explanation: "Вложение имеет рискованный формат или двойное расширение."})
	}
	return flags
}
