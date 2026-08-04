#!/usr/bin/env python3
"""
CV Matcher - CV Coach Dataset Generator (v2)

Gerçek canlı sistem promptu (buildSystemPrompt) formatını birebir taklit eden,
skor bağlamı gömülü, çok turlu sohbet örnekleri içeren kapsamlı eğitim verisi üretir.

Her örnekte:
  - system:  canlı sistem promptu yapısı (skor + CV + JD bağlamı)
  - input:   kullanıcı mesajı
  - output:  asistan cevabı (bağlamdaki skorlari AYNEN tekrar eder)
  - history: isteğe bağlı önceki turlar (çok turlu öğrenim)

Usage:
    python generate_cv_coach_dataset.py [--output data/cv_coach_dataset.json] [--count 800]
"""

import argparse
import json
import os
import random

random.seed(42)

# ============================================================
# SISTEM PROMPTLARI (canli buildSystemPrompt ile ayni yapi)
# ============================================================

SYSTEM_EN = (
    "You are the CV Coach, an expert career assistant. Help the user improve their CV, "
    "prepare for interviews, and tailor applications. Be concrete, practical and encouraging. "
    "Use short paragraphs and give examples with numbers when possible. "
    "ALWAYS answer in the same language the user writes in: if the user writes in Turkish, "
    "answer in Turkish (Türkçe); if they write in English, answer in English. "
    "Never switch languages mid-answer. When match results are provided below, quote the "
    "exact scores as written in the context (for example \"25/100\"). Never invent, estimate "
    "or recalculate scores — if a number is not in the context, say you do not have it."
)

SYSTEM_TR = (
    "Sen CV Koçusun, uzman kariyer asistanısın. Kullanıcının CV'sini iyileştirmesine, "
    "mülakatlara hazırlanmasına ve başvurularını uyarlamasına yardım et. Somut, pratik ve "
    "cesaretlendirici ol. Kısa paragraflar kullan ve mümkünse rakamlarla örnekler ver. "
    "HER ZAMAN kullanıcının yazdığı dilde cevap ver: kullanıcı Türkçe yazarsa Türkçe, "
    "İngilizce yazarsa İngilizce cevap ver. Cevabın ortasında asla dil değiştirme. "
    "Aşağıda eşleşme sonuçları verildiğinde, bağlamda yazan skorları aynen aktar "
    "(örneğin \"25/100\"). Skor uydurma, tahmin etme veya yeniden hesaplama — "
    "bir sayı bağlamda yoksa, elinizde olmadığını söyle."
)

# ============================================================
# SKOR PROFILLERI (öğretici bağlamlar)
# ============================================================

PROFILES = [
    {
        "cv_title": "Backend Developer",
        "cv_skills": ["Python", "Django", "PostgreSQL", "Docker"],
        "cv_summary": "Backend developer with 4 years of experience in Python and Django.",
        "jd_title": "Backend Engineer",
        "jd_skills": ["FastAPI", "Redis", "Kafka", "Kubernetes"],
        "overall": 25, "skill": 0, "exp": 50, "edu": 50,
        "matched": ["Python", "Docker"],
        "missing": ["FastAPI", "Redis", "Kafka", "Kubernetes"],
    },
    {
        "cv_title": "Backend Developer",
        "cv_skills": ["Python", "Django", "PostgreSQL", "Docker"],
        "cv_summary": "Backend developer with 4 years of experience in Python and Django.",
        "jd_title": "Backend Engineer",
        "jd_skills": ["Python", "Django", "Docker", "Redis"],
        "overall": 58, "skill": 62, "exp": 55, "edu": 50,
        "matched": ["Python", "Django", "Docker"],
        "missing": ["Redis"],
    },
    {
        "cv_title": "Full Stack Developer",
        "cv_skills": ["React", "Node.js", "MongoDB", "CSS"],
        "cv_summary": "Full stack developer passionate about building web apps.",
        "jd_title": "Senior Full Stack Engineer",
        "jd_skills": ["React", "Node.js", "TypeScript", "PostgreSQL"],
        "overall": 72, "skill": 80, "exp": 65, "edu": 60,
        "matched": ["React", "Node.js"],
        "missing": ["TypeScript", "PostgreSQL"],
    },
    {
        "cv_title": "Data Scientist",
        "cv_skills": ["Python", "Pandas", "scikit-learn", "SQL"],
        "cv_summary": "Data scientist with 5 years of experience in predictive modeling.",
        "jd_title": "ML Engineer",
        "jd_skills": ["PyTorch", "MLflow", "Docker", "Kubernetes"],
        "overall": 41, "skill": 30, "exp": 35, "edu": 70,
        "matched": ["Python", "SQL"],
        "missing": ["PyTorch", "MLflow", "Docker", "Kubernetes"],
    },
    {
        "cv_title": "DevOps Engineer",
        "cv_skills": ["AWS", "Docker", "CI/CD", "Linux"],
        "cv_summary": "DevOps engineer specializing in cloud infrastructure automation.",
        "jd_title": "DevOps Engineer",
        "jd_skills": ["Terraform", "Kubernetes", "ArgoCD", "Prometheus"],
        "overall": 47, "skill": 40, "exp": 60, "edu": 40,
        "matched": ["Docker", "Linux"],
        "missing": ["Terraform", "Kubernetes", "ArgoCD", "Prometheus"],
    },
    {
        "cv_title": "Frontend Developer",
        "cv_skills": ["React", "TypeScript", "Tailwind CSS", "Next.js"],
        "cv_summary": "Frontend developer who loves clean, accessible UI.",
        "jd_title": "Frontend Engineer",
        "jd_skills": ["React", "TypeScript", "Next.js", "Jest"],
        "overall": 85, "skill": 90, "exp": 80, "edu": 75,
        "matched": ["React", "TypeScript", "Next.js"],
        "missing": ["Jest"],
    },
    {
        "cv_title": "Mobile Developer",
        "cv_skills": ["Flutter", "Dart", "Firebase"],
        "cv_summary": "Mobile developer with experience building cross-platform apps.",
        "jd_title": "Mobile Developer",
        "jd_skills": ["Kotlin", "Swift", "React Native", "GraphQL"],
        "overall": 22, "skill": 10, "exp": 40, "edu": 55,
        "matched": ["Firebase"],
        "missing": ["Kotlin", "Swift", "React Native", "GraphQL"],
    },
    {
        "cv_title": "Backend Developer",
        "cv_skills": ["Java", "Spring Boot", "MySQL", "Kafka"],
        "cv_summary": "Java backend engineer with 6 years of experience in microservices.",
        "jd_title": "Backend Engineer",
        "jd_skills": ["Java", "Spring Boot", "Kafka", "Docker"],
        "overall": 68, "skill": 70, "exp": 65, "edu": 60,
        "matched": ["Java", "Spring Boot", "Kafka"],
        "missing": ["Docker"],
    },
    {
        "cv_title": "Full Stack Developer",
        "cv_skills": ["React", "Node.js", "MongoDB", "CSS"],
        "cv_summary": "Full stack developer passionate about building web apps.",
        "jd_title": "Backend Engineer",
        "jd_skills": ["Node.js", "TypeScript", "PostgreSQL", "Redis"],
        "overall": 55, "skill": 50, "exp": 60, "edu": 60,
        "matched": ["Node.js"],
        "missing": ["TypeScript", "PostgreSQL", "Redis"],
    },
    {
        "cv_title": "Data Scientist",
        "cv_skills": ["Python", "Pandas", "scikit-learn", "SQL"],
        "cv_summary": "Data scientist with 5 years of experience in predictive modeling.",
        "jd_title": "Data Engineer",
        "jd_skills": ["Spark", "Airflow", "dbt", "Snowflake"],
        "overall": 33, "skill": 25, "exp": 30, "edu": 60,
        "matched": ["Python", "SQL"],
        "missing": ["Spark", "Airflow", "dbt", "Snowflake"],
    },
]

# ============================================================
# SKOR SORULARI (bağlamdaki sayıları AYNEN tekrarlar)
# ============================================================

SCORE_QA_EN = [
    {
        "q": "What is my overall match score?",
        "a": "Your overall score is {overall}/100, with skill match at {skill}/100, experience at {exp}/100 and education at {edu}/100. The biggest gap is missing skills: {missing}. Focus on those first to raise your score.",
    },
    {
        "q": "Why is my skill match score only {skill}/100?",
        "a": "Because the job description asks for {jd_skills}, while your CV lists {cv_skills}. Only {matched} overlap, so the skill match is {skill}/100. The missing skills are {missing}. Add the ones you have actually used at work.",
    },
    {
        "q": "How can I improve my score from {overall}/100?",
        "a": "Your current score is {overall}/100. The weakest categories are skills ({skill}/100) and experience ({exp}/100). Missing skills from the job description: {missing}. Add these if you have used them, and rewrite your experience bullets with metrics to lift the experience score.",
    },
    {
        "q": "Is my score of {overall}/100 good?",
        "a": "{overall}/100 is {verdict_en}. The gap is mainly in {lowest_cat}: {missing}. With focused work on those skills and stronger bullet points, you can realistically reach 70/100 in a few weeks.",
    },
    {
        "q": "Which is higher: my skill match or my education score?",
        "a": "Your education score ({edu}/100) is higher than your skill match ({skill}/100). Your overall score is {overall}/100. Since skills weigh heavily in most roles, investing in the missing skills ({missing}) will move your overall score more than any other change.",
    },
    {
        "q": "What will raise my overall score the fastest?",
        "a": "Adding the missing skills ({missing}) to your CV — if you genuinely have experience with them. That alone can lift the skill match ({skill}/100) and the overall score. Second fastest: converting your experience bullets from responsibilities into measured results, which raises the experience score ({exp}/100).",
    },
    {
        "q": "How did you calculate my score?",
        "a": "Your overall score ({overall}/100) is composed of skill match ({skill}/100), experience ({exp}/100) and education ({edu}/100). The skill match compares your CV skills ({cv_skills}) against the job description requirements ({jd_skills}). Matched: {matched}. Missing: {missing}.",
    },
    {
        "q": "What happens if I add {missing0} to my CV?",
        "a": "If you genuinely have experience with {missing0}, adding it raises your skill match ({skill}/100) and therefore your overall score ({overall}/100). I would still recommend being honest — interviewers probe listed skills, especially during technical deep-dives.",
    },
    {
        "q": "My experience score is {exp}/100. How do I improve it?",
        "a": "Your experience score is {exp}/100. Rewrite every bullet with the STAR pattern and add one number per line: users, revenue, latency, team size or time saved. Instead of \"built REST APIs\", write \"shipped REST APIs serving 50K requests/day, cutting response time by 30%\". This is usually the fastest lever after fixing the missing skills ({missing}).",
    },
    {
        "q": "Can you explain my match breakdown?",
        "a": "Sure. Overall: {overall}/100. Skill match: {skill}/100 (matched: {matched}; missing: {missing}). Experience: {exp}/100. Education: {edu}/100. Your quickest win is addressing the missing skills and quantifying your experience bullets.",
    },
]

SCORE_QA_TR = [
    {
        "q": "Toplam eşleşme skorum kaç?",
        "a": "Toplam skorunuz {overall}/100: beceri eşleşmesi {skill}/100, deneyim {exp}/100, eğitim {edu}/100. En büyük eksik, becerilerde: {missing}. Skoru artırmak için önce bu eksik becerilere odaklanın.",
    },
    {
        "q": "Beceri eşleşme skorum neden sadece {skill}/100?",
        "a": "Çünkü iş ilanında {jd_skills} isteniyor, CV'nizde ise {cv_skills} listeleniyor. Yalnızca {matched} eşleşiyor, bu yüzden beceri eşleşmesi {skill}/100. Eksik beceriler: {missing}. İşte gerçekten kullandıklarınızı ekleyin.",
    },
    {
        "q": "{overall}/100'den skorumu nasıl yükseltebilirim?",
        "a": "Mevcut skorunuz {overall}/100. En zayıf kategoriler beceriler ({skill}/100) ve deneyim ({exp}/100). İlanda istenen ve eksik olan beceriler: {missing}. Bunları kullandıysanız ekleyin ve deneyim maddelerinizi ölçülerle yeniden yazın — bu deneyim skorunu yukarı çeker.",
    },
    {
        "q": "{overall}/100 iyi bir skor mu?",
        "a": "{overall}/100, {verdict_tr}. Ana boşluk {lowest_cat_tr}: {missing}. Eksik becerilere odaklanıp deneyim maddelerinizi güçlendirirseniz birkaç hafta içinde 70/100'e ulaşabilirsiniz.",
    },
    {
        "q": "Hangisi daha yüksek: beceri eşleşmem mi eğitim skorum mu?",
        "a": "Eğitim skorunuz ({edu}/100), beceri eşleşmenizden ({skill}/100) daha yüksek. Toplam skorunuz {overall}/100. Çoğu rolde beceriler daha ağırlıklı olduğu için, eksik becerileri ({missing}) eklemek toplam skoru en çok artıran değişikliktir.",
    },
    {
        "q": "Toplam skorumu en hızlı ne yükseltir?",
        "a": "Eksik becerileri ({missing}) CV'nize eklemek — tabii gerçekten deneyiminiz varsa. Bu tek başına beceri eşleşmesini ({skill}/100) ve toplam skoru yükseltir. İkinci en hızlı yol: deneyim maddelerinizi sorumluluk cümlelerinden ölçülü sonuçlara çevirmek; bu da deneyim skorunu ({exp}/100) artırır.",
    },
    {
        "q": "Skorumu nasıl hesapladınız?",
        "a": "Toplam skorunuz ({overall}/100) üç bileşenden oluşuyor: beceri eşleşmesi ({skill}/100), deneyim ({exp}/100) ve eğitim ({edu}/100). Beceri eşleşmesi, CV'nizdeki becerileri ({cv_skills}) ilandaki gereksinimlerle ({jd_skills}) karşılaştırır. Eşleşen: {matched}. Eksik: {missing}.",
    },
    {
        "q": "CV'me {missing0} eklersem ne olur?",
        "a": "{missing0} ile gerçekten deneyiminiz varsa eklemek, beceri eşleşmenizi ({skill}/100) ve toplam skorunuzu ({overall}/100) yükseltir. Yine de dürüst olmanızı öneririm — görüşmeciler listelenen becerileri derinlemesine sorgular.",
    },
    {
        "q": "Deneyim skorum {exp}/100. Nasıl iyileştiririm?",
        "a": "Deneyim skorunuz {exp}/100. Her maddeyi STAR yöntemiyle yeniden yazın ve satıra bir rakam ekleyin: kullanıcı, gelir, gecikme, ekip büyüklüğü veya kazanılan süre. \"REST API geliştirdim\" yerine \"günde 50K isteğe hizmet veren API'ler geliştirdim, yanıt süresini %30 kısalttım\" yazın. Eksik becerilerden ({missing}) sonra en hızlı kazanç budur.",
    },
    {
        "q": "Eşleşme dökümümü açıklar mısınız?",
        "a": "Tabii. Toplam: {overall}/100. Beceri eşleşmesi: {skill}/100 (eşleşen: {matched}; eksik: {missing}). Deneyim: {exp}/100. Eğitim: {edu}/100. En hızlı kazancınız eksik becerileri eklemek ve deneyim maddelerinizi ölçülendirmek.",
    },
]

# ============================================================
# GENEL KONU ŞABLONLARI (EN)
# ============================================================

QA_EN = [
    {
        "q": "What skills should I add to my CV for the {jd_title} role?",
        "a": "Compare your CV ({cv_skills}) to the {jd_title} job description. It lists {jd_skills} as requirements. If you have used these at work, add them to your skills section. If not, complete a short course and note them as \"familiar with\" rather than claiming expertise. Also mirror the posting's exact keywords — many companies use ATS filters that match resumes by exact wording.",
    },
    {
        "q": "How should I describe my experience as a {cv_title}?",
        "a": "Turn each responsibility into a STAR bullet: Situation, Task, Action, Result. Instead of \"built REST APIs\", write \"Built and shipped REST APIs serving 50K daily requests, cutting response time by 30%\". Start every bullet with a strong action verb (built, led, automated, reduced, designed) and add at least one metric per bullet: users, revenue, performance percentage, team size or time saved.",
    },
    {
        "q": "My CV summary says: \"{cv_summary}\" Is this good? How can I improve it?",
        "a": "It is a starting point, but it is generic. A strong summary has three parts: (1) who you are — \"{cv_title} with {years}+ years of experience\"; (2) your top differentiator; (3) a value statement for the employer. Aim for 2-3 sentences. Replace vague words like \"passionate\" and \"hard-working\" with specific skills, and mirror the {jd_title} posting keywords.",
    },
    {
        "q": "How do I tailor my CV for the {jd_title} role?",
        "a": "Step 1: extract the top 5-8 requirements from the posting ({jd_skills}). Step 2: for each requirement, match a line in your CV — add or reword bullets to include the keyword. Step 3: reorder skills so the most relevant appear first. Step 4: adjust your summary to echo the role title and stack. Step 5: remove unrelated bullets to keep it to one page. Aim for 80%+ keyword overlap.",
    },
    {
        "q": "I have an interview for the {jd_title} position next week. What should I prepare?",
        "a": "Be ready to (1) walk through every bullet on your CV with a short story — who, what, how, result; (2) explain the biggest technical challenge you faced and how you solved it; (3) answer why you want this specific company; (4) describe a conflict with a teammate; (5) give a salary expectation range. Also prepare 3 questions about the team, roadmap and success metrics.",
    },
    {
        "q": "Will an ATS read my CV correctly? Any formatting tips?",
        "a": "ATS systems struggle with tables, graphics, headers/footers and unusual fonts. Use: (1) standard section titles — Summary, Experience, Education, Skills; (2) a single-column layout; (3) standard fonts like Arial or Calibri; (4) spell out acronyms at least once; (5) .docx or PDF (not images). Verify by copying the text back out of the file.",
    },
    {
        "q": "What is the best structure for a one-page CV?",
        "a": "Top to bottom: (1) header with name, title, email, phone, LinkedIn; (2) 2-3 sentence summary; (3) core skills as comma-separated keywords; (4) experience, most recent first, 3-5 bullets per job with metrics; (5) education at the end; (6) optional certifications. Skip the \"References available on request\" line. Recruiters decide in about 6 seconds — make your fit obvious immediately.",
    },
    {
        "q": "I have a 1-year employment gap. How should I explain it?",
        "a": "Do not hide it — gaps are common. Keep dates consistent; if the gap included learning, add a brief line like \"2024: Completed AWS certification & built X project\". In interviews, have a 2-sentence explanation ready: what you did and what you learned. Frame it as a period of growth, not a void.",
    },
    {
        "q": "How do I add numbers to my CV if I don't have exact metrics?",
        "a": "Use relative estimates with qualifiers: \"served 10K+ users\", \"reduced response time by approx. 30%\", \"scaled to handle 3x traffic\", \"led a team of 4\", \"automated a task saving 2 hours/week\". If you cannot quantify, use scope: \"solely responsible for\", \"first engineer to build\", \"end-to-end ownership of\". Consistency matters more than precision.",
    },
    {
        "q": "Do I need a cover letter for the {jd_title} application?",
        "a": "Yes if the posting requests one or it is a small company. Keep it to 3 short paragraphs: (1) which role you are applying for and one sentence on why you fit; (2) your most relevant achievement with a number; (3) why this company, and a call to action. Do not repeat the CV line by line — add personality and context.",
    },
    {
        "q": "How should I organize my skills section?",
        "a": "Group skills into 3-4 categories: Languages/Frameworks, Tools, Cloud/DevOps, Soft Skills. Put the most job-relevant category first. Do not use rating bars or percentages — they are subjective and ATS-hostile. Limit to your strongest 15-20 skills. Match the exact terminology used in job descriptions.",
    },
    {
        "q": "I have a non-CS degree but work as a {cv_title}. How do I handle education?",
        "a": "List your degree as-is, then highlight relevant coursework or projects from it. Place a \"Certifications\" section right after education — above experience if it is more relevant than your degree. Emphasize continuous learning: mention recent courses. Recruiters care most about what you can do today.",
    },
    {
        "q": "Which certifications should I add as a {cv_title}?",
        "a": "Add any certification you hold in a dedicated section with the issuing body and year. For a {cv_title}, the most recognized ones are AWS Certified Developer and CKA. Only list certifications you can discuss in depth — interviewers will ask about them. Drop expired ones unless they are foundational.",
    },
    {
        "q": "Give me examples of strong bullet points for a {cv_title} CV.",
        "a": "Weak: \"Responsible for backend development\". Strong: \"Designed and shipped Python/Django APIs processing 50K requests/day, reducing average latency by 35%\". Weak: \"Worked with Docker\". Strong: \"Containerized 12 services with Docker Compose, cutting setup time from hours to minutes\". Weak: \"Good at teamwork\". Strong: \"Collaborated with 3 engineers and a PM to ship a feature used by 10K users within 2 sprints\".",
    },
    {
        "q": "What are the most common interview questions for a {cv_title}?",
        "a": "Typical questions: (1) \"Tell me about yourself\" — a 2-minute summary of your career arc; (2) \"Describe a project you're proud of\" — use STAR; (3) technical deep-dive on the tools in your CV — be honest about depth; (4) \"Where do you see yourself in 5 years?\" — align with the company's growth; (5) \"Why should we hire you?\" — summarize your unique combination of skills and experience.",
    },
]

# ============================================================
# GENEL KONU ŞABLONLARI (TR)
# ============================================================

QA_TR = [
    {
        "q": "{jd_title} rolü için CV'me hangi becerileri eklemeliyim?",
        "a": "CV'nizi ({cv_skills}) {jd_title} iş ilanıyla karşılaştırın. İlan {jd_skills} istiyor. Bunları işte kullandıysanız beceri bölümüne ekleyin; kullanmadıysanız kısa bir kurs alıp \"temel düzeyde biliyorum\" olarak belirtin — uzman gibi göstermeyin. İlanın anahtar kelimelerini birebir kullanın: birçok firma başvuruları ATS (Başvuru Takip Sistemi) ile tam kelime eşleşmesine göre filtreler.",
    },
    {
        "q": "{cv_title} olarak deneyimimi nasıl yazmalıyım?",
        "a": "Her sorumluluğu STAR yöntemiyle yazın: Durum, Görev, Eylem, Sonuç. \"REST API geliştirdim\" yerine \"Günde 50K isteğe hizmet veren REST API'ler geliştirdim, yanıt süresini %30 kısalttım\" yazın. Her maddeye güçlü bir eylem fiiliyle başlayın (geliştirdim, yönettim, otomatikleştirdim, azalttım, tasarladım) ve en az bir ölçü ekleyin: kullanıcı, gelir, performans yüzdesi, ekip büyüklüğü veya kazanılan süre.",
    },
    {
        "q": "CV özetim şu an: \"{cv_summary}\" İyi mi, nasıl geliştirebilirim?",
        "a": "Özetiniz doğru yolda ama daha güçlü olabilir. Formül: [unvan] + [yıl] deneyim + [en güçlü 2 beceri] + [somut bir başarı]. Örneğin: \"4 yıllık deneyimli Backend Developer; Python ve Django ile ölçeklenebilir API'ler geliştirdi, yanıt süresini %30 kısalttı.\" 2-3 cümleyle sınırlayın ve başvurduğunuz ilanın anahtar kelimelerini yerleştirin.",
    },
    {
        "q": "CV'mi {jd_title} ilanına nasıl uyarlarım?",
        "a": "Adım 1: ilandaki en önemli 5-8 gereksinimi çıkarın ({jd_skills}). Adım 2: her gereksinim için CV'nizde bir satır bulun — yoksa ekleyin veya mevcut maddeleri anahtar kelimeyi içerecek şekilde yazın. Adım 3: becerileri en alakalı olan başta olacak şekilde sıralayın. Adım 4: özeti rol başlığına göre güncelleyin. Adım 5: ilgisiz maddeleri çıkarıp tek sayfada tutun. %80+ anahtar kelime örtüşmesi hedefleyin.",
    },
    {
        "q": "Önümüzdeki hafta {jd_title} mülakatım var. Ne hazırlanmalıyım?",
        "a": "Şunlara hazır olun: (1) CV'nizdeki her maddeyi kısa bir hikayeyle anlatın — kim, ne, nasıl, sonuç; (2) yaşadığınız en büyük teknik zorluğu ve nasıl çözdüğünüzü anlatın; (3) neden bu şirket sorusuna cevap; (4) ekip arkadaşıyla yaşadığınız bir çatışmayı anlatın; (5) maaş beklentinizi söyleyin. Ayrıca ekip, yol haritası ve başarı ölçütleri hakkında 3 soru hazırlayın.",
    },
    {
        "q": "CV'mi ATS sistemlerinden nasıl geçirebilirim?",
        "a": "ATS (Başvuru Takip Sistemi) için: (1) standart başlıklar kullanın — \"Deneyim\", \"Eğitim\", \"Beceriler\"; (2) tablo, sütun ve resim içermeyen basit bir PDF düzeni seçin; (3) ilandaki anahtar kelimeleri birebir kullanın; (4) tarihleri \"Ocak 2022 - Haziran 2024\" gibi net yazın; (5) başlıklara veya görsellere metin gömmeyin. Kaydettikten sonra metni kopyalayıp doğrulayın.",
    },
    {
        "q": "Tek sayfalık bir CV'nin en iyi yapısı ne olmalı?",
        "a": "Yukarıdan aşağıya: (1) isim, unvan, e-posta, telefon, LinkedIn; (2) 2-3 cümlelik özet; (3) virgülle ayrılmış anahtar beceriler; (4) en yeni iş önce olacak şekilde deneyim — her iş için ölçülü 3-5 madde; (5) eğitim; (6) isteğe bağlı sertifikalar. \"Referanslar istek üzerine\" satırını koymayın. İşe alım uzmanı 6 saniyede karar verir — uygunluğunuz hemen görünmeli.",
    },
    {
        "q": "CV'mde 6 aylık bir iş deneyimi boşluğu var. Bunu nasıl açıklamalıyım?",
        "a": "Kısa boşlukları özetinize yazmanıza gerek yok, ama görüşmede hazırlıklı olun. Boşlukta ne yaptıysanız vurgulayın: kurs, freelance proje, gönüllülük, dil eğitimi veya kişisel proje. CV'ye \"Kariyer molası sırasında X eğitimini tamamladım\" gibi kısa bir satır ekleyebilirsiniz. Dürüst ama olumlu bir çerçeve çizin.",
    },
    {
        "q": "Başarılarım için rakam bulamıyorum. Ne yapmalıyım?",
        "a": "Rakamlar şu kaynaklardan gelebilir: işlem hacmi, kullanıcı sayısı, ekip boyutu, maliyet tasarrufu, kazanılan saat, azalan hata oranı. \"Destek sürecini otomatikleştirdim\" yerine \"otomatikleştirerek ayda 10 saat tasarruf sağladım\" yazın. Hiçbir veri yoksa git/rapor geçmişinden tahmin çıkarın ve \"yaklaşık\" notu düşün.",
    },
    {
        "q": "{jd_title} başvurusu için ön yazı yazmalı mıyım?",
        "a": "İlan istiyorsa veya şirket küçükse evet. 3 kısa paragraf tutun: (1) hangi pozisyona başvurduğunuz ve neden uygun olduğunuz tek cümle; (2) en alakalı başarınız ve rakamı; (3) neden bu şirket ve çağrı cümlesi. CV'nizi satır satır tekrar etmeyin — kişilik ve bağlam ekleyin. İlan istemiyorsa o vakti CV'nizi uyarlamaya harcayın.",
    },
    {
        "q": "Beceri bölümümü nasıl düzenlemeliyim?",
        "a": "Becerileri 3-4 gruba ayırın: Programlama Dilleri/Çerçeveler, Araçlar, Bulut/DevOps, Sosyal Beceriler. İşle en alakalı grubu en üste koyun. Puan çubuğu veya yüzde kullanmayın — subjektiftir ve ATS'ye uygun değildir. En güçlü 15-20 beceriyle sınırlayın; 40 beceri etkiyi azaltır. İlanlarda geçen terimleri birebir kullanın.",
    },
    {
        "q": "İlgisiz bir bölümden mezun oldum ama {cv_title} olarak çalışıyorum. Eğitimi nasıl yazmalıyım?",
        "a": "Derecenizi olduğu gibi listeleyin, ilgili dersleri ve projeleri öne çıkarın. Sertifikalarınız varsa eğitimden hemen sonra ayrı bir \"Sertifikalar\" bölümü açın — dereceden daha alakalıysa deneyimin üstüne bile koyabilirsiniz. Sürekli öğrenmeyi vurgulayın: yakın zamanda bitirdiğiniz kursları yazın.",
    },
    {
        "q": "{cv_title} olarak CV'me hangi sertifikaları eklemeliyim?",
        "a": "Sahip olduğunuz tüm sertifikaları veren kurum ve yılıyla birlikte ayrı bir bölüme ekleyin. {cv_title} için en çok tanınanlar AWS Certified Developer ve CKA'dır. Sadece görüşmede derinlemesine anlatabileceğiniz sertifikaları yazın — görüşmeciler mutlaka soracaktır. Temel olanlar dışında süresi geçmişleri çıkarın.",
    },
    {
        "q": "{cv_title} CV'si için güçlü madde örnekleri verir misiniz?",
        "a": "Zayıf: \"Backend geliştirmeden sorumluydum\". Güçlü: \"Günde 50K isteğe hizmet veren Python/Django API'ler geliştirdim, ortalama gecikmeyi %35 azalttım\". Zayıf: \"Docker kullandım\". Güçlü: \"Docker Compose ile 12 servisi konteynerize ettim, kurulum süresini saatlerden dakikalara indirdim\". Zayıf: \"Ekip çalışmasına yatkınım\". Güçlü: \"2 sprint içinde 10K kullanıcının kullandığı özelliği 3 mühendis ve bir PM ile teslim ettik\".",
    },
    {
        "q": "{cv_title} için en sık sorulan mülakat soruları neler?",
        "a": "Tipik sorular: (1) \"Kendini tanıt\" — kariyer yolculuğunuzu 2 dakikada özetleyin; (2) \"Gurur duyduğun bir projeyi anlat\" — STAR kullanın; (3) CV'nizdeki araçlarla ilgili teknik derinlemesine soru — derinlik konusunda dürüst olun; (4) \"5 yıl sonra kendini nerede görüyorsun?\" — şirketin büyümesiyle uyumlu olun; (5) \"Seni neden işe alalım?\" — beceri ve deneyim kombinasyonunuzu özetleyin.",
    },
]

# ============================================================
# ÇOK TURLU DEVAM SORULARI
# ============================================================

FOLLOWUP_EN = [
    ("What about my education score?", "Your education score is {edu}/100. It is {edu_comment}. I would leave it as-is and focus on skills ({skill}/100) and experience ({exp}/100), since those have the most room to grow."),
    ("Should I add {missing0} to my CV?", "Only if you have real experience with {missing0} — interviewers will probe it. If you do, add it and it will directly raise your skill match ({skill}/100) and overall score ({overall}/100)."),
    ("What is the fastest improvement I can make this week?", "This week: add the missing skills you genuinely know ({missing0}) and rewrite your top 2 experience bullets with numbers. Those two changes alone should visibly raise your {overall}/100 score."),
    ("How long will it take to reach 70/100?", "Depending on the gaps: adding real skills can take weeks or months; polishing bullets and keywords can take days. With {missing} missing, a realistic short-term target is 55-60/100, then 70+ once the skills are genuinely added."),
]

FOLLOWUP_TR = [
    ("Peki eğitim skorum ne durumda?", "Eğitim skorunuz {edu}/100. Bu oldukça iyi. Önceliğinizi becerilere ({skill}/100) ve deneyime ({exp}/100) verin — büyüme alanı en çok onlarda."),
    ("CV'me {missing0} eklemeli miyim?", "Sadece gerçek deneyiminiz varsa ekleyin — görüşmeciler sorgular. Varsa ekleyin, beceri eşleşmenizi ({skill}/100) ve toplam skorunuzu ({overall}/100) doğrudan yükseltir."),
    ("Bu hafta yapabileceğim en hızlı iyileştirme ne?", "Bu hafta: gerçekten bildiğiniz eksik becerileri ({missing0}) ekleyin ve deneyiminizdeki en iyi 2 maddeyi rakamlarla yeniden yazın. İkisi birlikte {overall}/100 skorunuzu gözle görülür biçimde yükseltir."),
    ("70/100'e ulaşmak ne kadar sürer?", "Boşluğa bağlı: gerçek beceri eklemek haftalar veya aylar sürebilir; maddeleri ve anahtar kelimeleri cilalamak günler. {missing} eksikken gerçekçi kısa vade hedefi 55-60/100, beceriler gerçekten eklendikten sonra 70+."),
]

# ============================================================
# YARDIMCILAR
# ============================================================

def build_system(profile, lang, include_score=True):
    base = SYSTEM_TR if lang == "tr" else SYSTEM_EN
    parts = [base]
    if include_score:
        parts.append(
            f"\n\nThe AI match result between the user's CV and this job description "
            f"(the user wants to understand and improve this score):\n"
            f"Overall score: {profile['overall']}/100\n"
            f"Skill match: {profile['skill']}/100, Experience: {profile['exp']}/100, Education: {profile['edu']}/100\n"
            f"Matched skills: {', '.join(profile['matched'])}\n"
            f"Missing skills: {', '.join(profile['missing'])}"
        )
    parts.append(
        f"\n\nThe user's CV context (use it to give personalized advice):\n"
        f"Title: {profile['cv_title']}\n"
        f"Skills: {', '.join(profile['cv_skills'])}\n"
        f"Summary: {profile['cv_summary']}"
    )
    parts.append(
        f"\n\nThe target job description (the user is applying to this role):\n"
        f"Title: {profile['jd_title']}\n"
        f"Required skills: {', '.join(profile['jd_skills'])}"
    )
    return "".join(parts)


def fmt(template, profile):
    return template.format(
        overall=profile["overall"],
        skill=profile["skill"],
        exp=profile["exp"],
        edu=profile["edu"],
        matched=", ".join(profile["matched"]),
        missing=", ".join(profile["missing"]),
        missing0=profile["missing"][0] if profile["missing"] else "the missing skills",
        cv_title=profile["cv_title"],
        cv_skills=", ".join(profile["cv_skills"]),
        cv_summary=profile["cv_summary"],
        jd_title=profile["jd_title"],
        jd_skills=", ".join(profile["jd_skills"]),
        years=random.randint(2, 8),
        lowest_cat="skills" if profile["skill"] <= min(profile["exp"], profile["edu"]) else "experience",
        lowest_cat_tr="beceriler" if profile["skill"] <= min(profile["exp"], profile["edu"]) else "deneyim",
        edu_comment="very solid" if profile["edu"] >= 60 else "fine",
        verdict_en="a solid score" if profile["overall"] >= 70 else ("a decent starting point" if profile["overall"] >= 50 else "below average for most postings"),
        verdict_tr="çok iyi bir skor" if profile["overall"] >= 70 else ("iyi bir başlangıç" if profile["overall"] >= 50 else "çoğu ilan için ortalamanın altında"),
    )


def generate(count, score_ratio=0.45, multi_turn_ratio=0.25):
    items = []
    i = 0
    while len(items) < count:
        profile = random.choice(PROFILES)
        lang = "tr" if i % 2 == 1 else "en"
        is_score = (i % 2 == 1) or (random.random() < score_ratio) if i > 2 else (i % 2 == 0)
        is_score = (i % 4 < 2)  # iki dilde de bol skor sorusu
        i += 1

        if is_score:
            qa_pool = SCORE_QA_TR if lang == "tr" else SCORE_QA_EN
            qa = random.choice(qa_pool)
            question = fmt(qa["q"], profile)
            answer = fmt(qa["a"], profile)
        else:
            qa_pool = QA_TR if lang == "tr" else QA_EN
            qa = random.choice(qa_pool)
            question = fmt(qa["q"], profile)
            answer = fmt(qa["a"], profile)

        system = build_system(profile, lang)

        item = {
            "system": system,
            "input": question,
            "output": answer,
        }

        # Çok turlu devam örnekleri
        if random.random() < multi_turn_ratio:
            follow_pool = FOLLOWUP_TR if lang == "tr" else FOLLOWUP_EN
            fq, fa = random.choice(follow_pool)
            item["history"] = [
                {"role": "user", "content": question},
                {"role": "assistant", "content": answer},
            ]
            item["input"] = fmt(fq, profile)
            item["output"] = fmt(fa, profile)

        items.append(item)

    return items


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--output", default="data/cv_coach_dataset.json")
    parser.add_argument("--count", type=int, default=800)
    args = parser.parse_args()

    items = generate(args.count)

    os.makedirs(os.path.dirname(args.output) or ".", exist_ok=True)
    with open(args.output, "w", encoding="utf-8") as f:
        json.dump(items, f, ensure_ascii=False, indent=2)

    tr = sum(1 for x in items if x["system"].startswith("Sen CV Koçusun"))
    multi = sum(1 for x in items if "history" in x)
    score = sum(1 for x in items if any(k in x["input"] for k in ["skor", "score", "eşleşme", "match", "döküm", "breakdown"]))
    print(f"[OK] {len(items)} örnek {args.output} dosyasına yazıldı")
    print(f"  TR örnek: {tr}, EN örnek: {len(items)-tr}, çok turlu: {multi}, skor sorusu: {score}")


if __name__ == "__main__":
    main()
