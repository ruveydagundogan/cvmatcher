#!/usr/bin/env python3
"""
CV Matcher - CV Coach Dataset Generator

CV geliştirme asistanı için sohbet (Q&A) eğitim verisi üretir.
Konular: eksik skill tespiti, deneyim yazımı, özet iyileştirme,
JD'ye göre uyarlama, mülakat hazırlığı, ATS/format önerileri.

Usage:
    python generate_cv_coach_dataset.py [--output data/cv_coach_dataset.json] [--count 80]
"""

import argparse
import json
import random
import os

random.seed(7)

COACH_INSTRUCTIONS = [
    "You are a CV Coach. Help the user improve their CV with concrete, actionable advice.",
    "You are a friendly career coach who specializes in resume writing. Give specific, practical suggestions.",
    "You are an experienced HR professional and CV consultant. Help the user strengthen their CV.",
    "You are a CV coaching assistant. Be encouraging but honest, and always give examples.",
]

# ============================================================
# PERSONAS (CV örnekleri)
# ============================================================

PERSONAS = [
    {
        "role": "Backend Developer",
        "skills": ["Python", "Django", "PostgreSQL", "Docker"],
        "years": 4,
        "company": "TechCorp",
        "achievement": "built REST APIs",
        "summary": "Backend developer with 4 years of experience in Python and Django.",
    },
    {
        "role": "Full Stack Developer",
        "skills": ["React", "Node.js", "MongoDB", "CSS"],
        "years": 3,
        "company": "StartupX",
        "achievement": "developed user interfaces",
        "summary": "Full stack developer passionate about building web apps.",
    },
    {
        "role": "Data Scientist",
        "skills": ["Python", "Pandas", "scikit-learn", "SQL"],
        "years": 5,
        "company": "FinTechLabs",
        "achievement": "built ML models",
        "summary": "Data scientist with 5 years of experience in predictive modeling.",
    },
    {
        "role": "DevOps Engineer",
        "skills": ["AWS", "Docker", "CI/CD", "Linux"],
        "years": 6,
        "company": "CloudInc",
        "achievement": "managed cloud infrastructure",
        "summary": "DevOps engineer specializing in cloud infrastructure automation.",
    },
    {
        "role": "Frontend Developer",
        "skills": ["React", "TypeScript", "Tailwind CSS", "Next.js"],
        "years": 2,
        "company": "DesignHub",
        "achievement": "built responsive interfaces",
        "summary": "Frontend developer who loves clean, accessible UI.",
    },
    {
        "role": "Mobile Developer",
        "skills": ["Flutter", "Dart", "Firebase"],
        "years": 3,
        "company": "AppWorks",
        "achievement": "shipped mobile apps",
        "summary": "Mobile developer with experience building cross-platform apps.",
    },
]

# ============================================================
# TEMPLATES
# ============================================================

QA_TEMPLATES = [
    # --- Eksik skill / JD uyarlama ---
    {
        "topic": "missing_skills",
        "template": "My CV says: {summary} My skills are: {skills}. I want to apply for a {target_role} position. What skills should I add to my CV?",
        "answers": [
            "Compare your CV to the {target_role} job description. Most {target_role} roles list {target_skills} as requirements. If you have used these at work, add them to your skills section. If not, consider completing a short course and noting it as \"familiar with\" rather than claiming expertise. Also re-read the job description line by line and mirror its keywords — many companies use ATS (Applicant Tracking Systems) that filter resumes by exact keyword matches.",
        ],
    },
    {
        "topic": "missing_skills",
        "template": "I'm a {role} with {years} years of experience. My CV lists {skills}. What's missing that recruiters look for?",
        "answers": [
            "For a {role} with {years} years of experience, recruiters typically look for: (1) measurable achievements, not just responsibilities — add numbers like \"reduced latency by 40%\"; (2) the tools most used in your niche — for {role}, tools like {target_skills} are commonly listed in postings; (3) soft skills demonstrated with examples, e.g. \"led a team of 3\"; (4) a clear summary sentence stating your years of experience and top 3 strengths.",
        ],
    },
    # --- Deneyim yazımı ---
    {
        "topic": "experience",
        "template": "How should I describe my experience as a {role} at {company}? I {achievement} and that's about it.",
        "answers": [
            "Turn each responsibility into a STAR bullet: Situation, Task, Action, Result. Instead of \"{achievement}\", write \"Built and shipped REST APIs serving 50K daily requests, cutting response time by 30%\". Start every bullet with a strong action verb (built, led, automated, reduced, designed). Add at least one metric per bullet: users, revenue, performance percentage, team size, or time saved. If you don't know exact numbers, estimate and say \"approx.\".",
        ],
    },
    {
        "topic": "experience",
        "template": "My experience section is just one line per job: \"{achievement}\". How do I make it stronger?",
        "answers": [
            "Expand each job into 3-5 bullet points using this structure: (1) a one-line scope statement (\"Responsible for the API layer of a fintech platform\"), (2) 2-3 action bullets with metrics (\"Reduced p99 latency from 800ms to 200ms by optimizing queries\"), (3) one impact bullet (\"Enabling 10K daily active users\"). Use present tense for your current job and past tense for previous ones. Keep bullets under two lines each.",
        ],
    },
    # --- Özet iyileştirme ---
    {
        "topic": "summary",
        "template": "My CV summary currently says: \"{summary}\" Is this good? How can I improve it?",
        "answers": [
            "That summary is a starting point but it's generic. A strong summary has three parts: (1) who you are: \"{role} with {years}+ years of experience\"; (2) your top differentiator: e.g. \"specializing in high-traffic web applications\"; (3) a value statement for the employer: \"seeking to bring scalable architecture skills to a growing product team\". Aim for 2-3 sentences max. Replace vague words like \"passionate\" and \"hard-working\" with specific skills.",
        ],
    },
    {
        "topic": "summary",
        "template": "How long should my CV summary be and what should it include?",
        "answers": [
            "Keep it to 2-3 sentences. Include: your title and years of experience, your top 2-3 hard skills, and one achievement or value proposition. Do NOT mention hobbies, marital status, or nationality. Write it after you finish the rest of the CV so it summarizes your actual strengths. Customize the last sentence per application.",
        ],
    },
    # --- JD'ye göre uyarlama ---
    {
        "topic": "tailoring",
        "template": "I'm applying to a {target_role} job. How do I tailor my generic CV?",
        "answers": [
            "Step 1: extract the top 5-8 requirements from the job description. Step 2: for each requirement, find the matching line in your CV — if it's missing, add it or reword existing bullets to include the keyword. Step 3: reorder your skills so the most relevant ones appear first. Step 4: adjust your summary to echo the role title and key stack. Step 5: remove unrelated bullets to keep the CV to one page. Target a 80%+ keyword overlap with the posting.",
        ],
    },
    {
        "topic": "tailoring",
        "template": "The job description requires {target_skills} but I only know {skills}. Should I apply?",
        "answers": [
            "Yes — most job descriptions are wish lists. Apply if you meet 60-70% of the requirements. In your CV, make your transferable skills visible: e.g. if the job wants {target_skills} and you know {skills}, emphasize the shared concepts and mention that you're actively learning the gap. During the interview, show a project where you quickly ramped up on a new technology.",
        ],
    },
    # --- Mülakat hazırlığı ---
    {
        "topic": "interview",
        "template": "I have an interview for a {target_role} position next week. What should I prepare based on my CV?",
        "answers": [
            "Be ready to (1) walk through every bullet on your CV with a short story — who, what, how, result; (2) explain the biggest technical challenge you faced as a {role} and how you solved it; (3) answer why you want this specific company; (4) describe a conflict with a teammate; (5) give your salary expectation range. Also prepare 3 questions to ask the interviewer about the team, roadmap, and success metrics.",
        ],
    },
    {
        "topic": "interview",
        "template": "What are the most common interview questions for a {role}?",
        "answers": [
            "Typical questions: (1) \"Tell me about yourself\" — answer with a 2-minute summary of your career arc; (2) \"Describe a project you're proud of\" — use STAR; (3) technical deep-dive on the tools in your CV — be honest about depth; (4) \"Where do you see yourself in 5 years?\" — align with the company's growth; (5) \"Why should we hire you?\" — summarize your unique combination of skills and experience.",
        ],
    },
    # --- Yapı / format / ATS ---
    {
        "topic": "ats",
        "template": "Will an ATS read my CV correctly? Any formatting tips?",
        "answers": [
            "ATS systems struggle with tables, graphics, headers/footers, and unusual fonts. Use: (1) standard section titles — Summary, Experience, Education, Skills; (2) a single-column layout; (3) standard fonts like Arial or Calibri; (4) spell out acronyms at least once (e.g. \"CI/CD (Continuous Integration/Continuous Deployment)\"); (5) .docx or PDF (not images). Save as a plain, clean layout and verify by copying the text out of the file.",
        ],
    },
    {
        "topic": "structure",
        "template": "What's the best structure for a one-page CV?",
        "answers": [
            "Top to bottom: (1) header with name, title, email, phone, LinkedIn; (2) 2-3 sentence summary; (3) core skills as comma-separated keywords; (4) experience, most recent first, 3-5 bullets per job with metrics; (5) education at the end; (6) optional certifications. Skip the \"References available on request\" line. Keep everything that doesn't add value out — the goal is a recruiter can find your fit in 6 seconds.",
        ],
    },
    # --- Eğitim ---
    {
        "topic": "education",
        "template": "I have a non-CS degree but I work as a {role}. How should I handle the education section?",
        "answers": [
            "List your degree as-is, but highlight the relevant coursework or projects from it. If you have certifications or bootcamps (e.g. AWS, Kubernetes, a coding bootcamp), place a \"Certifications\" section right after education and put it ABOVE experience if it's more relevant than your degree. Emphasize continuous learning: mention recent courses you completed. Recruiters care most about what you can do today.",
        ],
    },
    {
        "topic": "education",
        "template": "Should I include my GPA and graduation year in my CV?",
        "answers": [
            "Include GPA only if it's 3.0+ and you have less than 3 years of experience. Include the graduation year if you graduated within the last 5 years; after that, omit it to avoid age bias. List the most recent degree first. If you have multiple degrees, group them under one \"Education\" section.",
        ],
    },
    # --- Kariyer boşluğu ---
    {
        "topic": "gap",
        "template": "I have a 1-year employment gap in my CV. How should I explain it?",
        "answers": [
            "Don't hide it — gaps are common. In your CV, you don't need a section for it; just keep dates consistent. If the gap included learning (courses, certificates, freelance, personal projects), add a brief line like \"2024: Completed AWS certification & built X project\". In interviews, have a 2-sentence explanation ready: what you did and what you learned. Frame it as a period of growth, not a void.",
        ],
    },
    # --- Ölçülebilir başarı ---
    {
        "topic": "metrics",
        "template": "How do I add numbers to my CV if I don't have exact metrics?",
        "answers": [
            "Use relative estimates with qualifiers: \"served 10K+ users\", \"reduced response time by approx. 30%\", \"scaled to handle 3x traffic\", \"led a team of 4\", \"automated a task saving 2 hours/week\". If you genuinely can't quantify, use scope: \"solely responsible for\", \"first engineer to build\", \"end-to-end ownership of\". Consistency matters more than precision.",
        ],
    },
    {
        "topic": "metrics",
        "template": "Give me examples of strong bullet points for a {role} CV.",
        "answers": [
            "Weak: \"Responsible for backend development\". Strong: \"Designed and shipped Python/Django APIs processing 50K requests/day, reducing average latency by 35%\". Weak: \"Worked with Docker\". Strong: \"Containerized 12 services with Docker Compose, cutting setup time from hours to minutes\". Weak: \"Good at teamwork\". Strong: \"Collaborated with 3 engineers and a PM to ship a feature used by 10K users within 2 sprints\".",
        ],
    },
    # --- Kapak mektubu ---
    {
        "topic": "cover_letter",
        "template": "Do I need a cover letter for my {target_role} application?",
        "answers": [
            "Yes if the posting requests one or it's a small company. Keep it to 3 short paragraphs: (1) which role you're applying for and one sentence on why you're a fit; (2) your most relevant achievement with a number; (3) why this company, and a call to action. Don't repeat the CV line by line — add personality and context. If the posting doesn't ask for one, spend that time tailoring the CV instead.",
        ],
    },
    # --- Sertifikalar ---
    {
        "topic": "certifications",
        "template": "Which certifications should I add to my CV as a {role}?",
        "answers": [
            "Add any certification you have — put them in a dedicated section with the issuing body and year. For a {role}, the most recognized ones are: {certs}. Only list certifications you can discuss in depth; interviewers will ask about them. Keep expired ones only if they're foundational (e.g. security clearances).",
        ],
    },
    # --- Beceri bölümü ---
    {
        "topic": "skills",
        "template": "How should I organize my skills section?",
        "answers": [
            "Group skills into 3-4 categories: Languages/Frameworks, Tools, Cloud/DevOps, Soft Skills. Put the most job-relevant category first. Don't use rating bars or percentages — they're subjective and ATS-hostile. Limit the section to your strongest 15-20 skills; listing 40 dilutes impact. Match the exact terminology used in job descriptions.",
        ],
    },
]

# ============================================================
# Veri üretimi
# ============================================================

TARGET_ROLES = [
    {"role": "Backend Engineer", "skills": "FastAPI, Redis, Kafka, Kubernetes", "certs": "AWS Certified Developer, CKA"},
    {"role": "Senior Frontend Engineer", "skills": "Redux, GraphQL, Next.js, Jest", "certs": "AWS Certified Cloud Practitioner, Meta Frontend"},
    {"role": "Data Engineer", "skills": "Spark, Airflow, dbt, Snowflake", "certs": "GCP Professional Data Engineer, Databricks"},
    {"role": "DevOps Engineer", "skills": "Terraform, Kubernetes, ArgoCD, Prometheus", "certs": "CKA, AWS DevOps Engineer"},
    {"role": "ML Engineer", "skills": "PyTorch, MLflow, Docker, ONNX", "certs": "AWS ML Specialty, TensorFlow Developer"},
    {"role": "Mobile Developer", "skills": "Kotlin, Swift, React Native, GraphQL", "certs": "Google Associate Android Developer"},
]


def generate(count):
    items = []
    for i in range(count):
        persona = random.choice(PERSONAS)
        target = random.choice(TARGET_ROLES)
        qa = random.choice(QA_TEMPLATES)
        instruction = random.choice(COACH_INSTRUCTIONS)
        answer_template = random.choice(qa["answers"])

        question = qa["template"].format(
            role=persona["role"],
            years=persona["years"],
            company=persona["company"],
            achievement=persona["achievement"],
            summary=persona["summary"],
            skills=", ".join(persona["skills"]),
            target_role=target["role"],
            target_skills=target["skills"],
            certs=target["certs"],
        )

        answer = answer_template.format(
            role=persona["role"],
            years=persona["years"],
            company=persona["company"],
            achievement=persona["achievement"],
            summary=persona["summary"],
            skills=", ".join(persona["skills"]),
            target_role=target["role"],
            target_skills=target["skills"],
            certs=target["certs"],
        )

        items.append({
            "instruction": instruction,
            "input": question,
            "output": answer,
        })

    return items


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--output", default="data/cv_coach_dataset.json")
    parser.add_argument("--count", type=int, default=80)
    args = parser.parse_args()

    items = generate(args.count)

    os.makedirs(os.path.dirname(args.output) or ".", exist_ok=True)
    with open(args.output, "w", encoding="utf-8") as f:
        json.dump(items, f, ensure_ascii=False, indent=2)

    print(f"[OK] {len(items)} örnek {args.output} dosyasına yazıldı")


if __name__ == "__main__":
    main()
