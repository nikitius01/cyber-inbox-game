package ai

const Prompt = `Ты генерируешь безопасные обучающие письма для игры "Инспектор входящих".
Верни только валидный JSON без Markdown.
Схема: id, category, senderName, senderEmail, subject, body, isPhishing, links, attachments, raw, redFlags, explanation.
links должны содержать visibleText, actualUrl, domain, protocol, isShortened, isSuspicious, riskReason.
attachments должны содержать fileName, displayName, extension, mimeType, sizeKb, hash, isSuspicious, riskReason.
raw должен содержать headers, bodyText, bodyHtml, source.
Не создавай реальный вредоносный код, рабочие payload, инструкции атаки или настоящие секреты.
Текст письма должен быть на русском языке.`
