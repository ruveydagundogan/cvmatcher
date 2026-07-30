#!/usr/bin/env python3
"""
CV Matcher Dataset Generator
Gerçekçi CV ve JD örnekleri üretir, fine-tuning için hazırlar.
"""

import json
import random
import os

random.seed(42)

# ============================================================
# TEMPLATES - Gerçekçi CV ve JD şablonları
# ============================================================

ROLES = [
    {
        "title": "Backend Go Developer",
        "title_tr": "Backend Go Geliştirici",
        "skills": ["Go", "PostgreSQL", "Docker", "Kubernetes", "Redis", "gRPC", "Kafka", "AWS", "Microservices", "CI/CD"],
        "experience_years": 5,
        "education": {"degree": "Bachelor of Science in Computer Engineering", "institution": "Istanbul Technical University"},
        "projects": ["Built a distributed payment processing system handling 10K req/s",
                     "Designed microservices architecture for an e-commerce platform",
                     "Implemented CI/CD pipelines reducing deployment time by 80%"],
    },
    {
        "title": "Senior Full Stack Developer",
        "title_tr": "Kıdemli Full Stack Geliştirici",
        "skills": ["React", "TypeScript", "Node.js", "PostgreSQL", "Docker", "GraphQL", "Next.js", "Redis", "AWS", "Tailwind CSS"],
        "experience_years": 7,
        "education": {"degree": "Master of Science in Computer Science", "institution": "Bogazici University"},
        "projects": ["Developed a real-time collaboration platform with 50K+ users",
                     "Led migration from monolith to microservices architecture",
                     "Built design system used by 5 product teams"],
    },
    {
        "title": "Data Scientist",
        "title_tr": "Veri Bilimci",
        "skills": ["Python", "TensorFlow", "PyTorch", "SQL", "Spark", "MLflow", "Docker", "AWS SageMaker", "NLP", "Computer Vision"],
        "experience_years": 4,
        "education": {"degree": "PhD in Machine Learning", "institution": "METU"},
        "projects": ["Developed NLP model for Turkish text classification with 94% accuracy",
                     "Built recommendation system improving engagement by 35%",
                     "Created MLOps pipeline for model deployment and monitoring"],
    },
    {
        "title": "DevOps Engineer",
        "title_tr": "DevOps Mühendisi",
        "skills": ["Docker", "Kubernetes", "Terraform", "Ansible", "AWS", "Linux", "CI/CD", "Prometheus", "Grafana", "Python"],
        "experience_years": 6,
        "education": {"degree": "Bachelor of Science in Computer Engineering", "institution": "Hacettepe University"},
        "projects": ["Managed 200+ microservices on Kubernetes with 99.99% uptime",
                     "Automated infrastructure provisioning reducing costs by 40%",
                     "Implemented monitoring stack serving 100M+ data points daily"],
    },
    {
        "title": "Frontend Developer",
        "title_tr": "Frontend Geliştirici",
        "skills": ["React", "TypeScript", "JavaScript", "CSS", "HTML", "Next.js", "Vue.js", "Webpack", "Jest", "Storybook"],
        "experience_years": 3,
        "education": {"degree": "Bachelor of Science in Computer Science", "institution": "Yildiz Technical University"},
        "projects": ["Built responsive web app used by 100K+ daily active users",
                     "Migrated legacy jQuery codebase to React, improving performance by 60%",
                     "Created component library with 50+ reusable UI components"],
    },
    {
        "title": "Machine Learning Engineer",
        "title_tr": "Makine Öğrenmesi Mühendisi",
        "skills": ["Python", "TensorFlow", "PyTorch", "Docker", "Kubernetes", "AWS", "MLflow", "SQL", "ONNX", "CUDA"],
        "experience_years": 4,
        "education": {"degree": "Master of Science in Artificial Intelligence", "institution": "KUIS AI Center"},
        "projects": ["Deployed real-time inference pipeline serving 100K requests/hour",
                     "Optimized model inference latency from 200ms to 15ms using TensorRT",
                     "Built automated ML pipeline reducing experiment cycle from weeks to hours"],
    },
    {
        "title": "QA Engineer",
        "title_tr": "Test Mühendisi",
        "skills": ["Selenium", "Cypress", "JUnit", "TestNG", "Python", "Java", "Docker", "Jenkins", "Postman", "JMeter"],
        "experience_years": 4,
        "education": {"degree": "Bachelor of Science in Software Engineering", "institution": "Gebze Technical University"},
        "projects": ["Built automated test suite covering 95% of critical paths",
                     "Reduced regression testing time from 3 days to 2 hours",
                     "Implemented performance testing identifying bottlenecks reducing response time by 50%"],
    },
    {
        "title": "Product Manager",
        "title_tr": "Ürün Yöneticisi",
        "skills": ["Product Strategy", "Roadmapping", "Agile", "Scrum", "Data Analysis", "A/B Testing", "SQL", "Figma", "JIRA", "User Research"],
        "experience_years": 6,
        "education": {"degree": "Master of Business Administration", "institution": "Sabanci University"},
        "projects": ["Launched SaaS product generating $2M ARR in first year",
                     "Led cross-functional team of 12 engineers and designers",
                     "Defined product strategy that increased user retention by 45%"],
    },
    {
        "title": "Security Engineer",
        "title_tr": "Güvenlik Mühendisi",
        "skills": ["Python", "Go", "Kubernetes", "Docker", "AWS Security", "Penetration Testing", "SIEM", "Cryptography", "OWASP", "Zero Trust"],
        "experience_years": 5,
        "education": {"degree": "Bachelor of Science in Electrical & Electronics Engineering", "institution": "Bilkent University"},
        "projects": ["Led SOC 2 Type II certification achieving zero findings",
                     "Implemented zero-trust architecture securing 5000+ endpoints",
                     "Developed automated security scanning pipeline preventing 200+ vulnerabilities"],
    },
    {
        "title": "iOS Developer",
        "title_tr": "iOS Geliştirici",
        "skills": ["Swift", "UIKit", "SwiftUI", "Core Data", "Firebase", "Xcode", "Combine", "RxSwift", "App Store Connect", "Fastlane"],
        "experience_years": 4,
        "education": {"degree": "Bachelor of Science in Computer Engineering", "institution": "Istanbul University"},
        "projects": ["Published 3 apps with 500K+ total downloads on App Store",
                     "Built offline-first architecture reducing API calls by 70%",
                     "Implemented CI/CD pipeline with automated TestFlight distribution"],
    },
]

JD_TEMPLATES = [
    {
        "title": "Senior Backend Developer",
        "company": "Trendyol",
        "description": "We are looking for a Senior Backend Developer to join our growing engineering team.",
        "required_skills": ["Go", "PostgreSQL", "Docker", "Kubernetes", "Microservices"],
        "preferred_skills": ["Kafka", "Redis", "gRPC", "AWS"],
        "experience": "5+ years",
        "type": "Full-time",
    },
    {
        "title": "Full Stack Developer",
        "company": "Getir",
        "description": "Join our product team to build next-generation delivery platform.",
        "required_skills": ["React", "TypeScript", "Node.js", "PostgreSQL", "Docker"],
        "preferred_skills": ["GraphQL", "Redis", "AWS", "Next.js"],
        "experience": "3+ years",
        "type": "Full-time",
    },
    {
        "title": "Data Scientist",
        "company": "Peakup",
        "description": "Help us solve complex business problems using machine learning and AI.",
        "required_skills": ["Python", "Machine Learning", "SQL", "Deep Learning", "NLP"],
        "preferred_skills": ["TensorFlow", "PyTorch", "Docker", "Spark"],
        "experience": "3+ years",
        "type": "Full-time",
    },
    {
        "title": "DevOps Lead",
        "company": "Hepsiburada",
        "description": "Lead our infrastructure team to build scalable cloud solutions.",
        "required_skills": ["Kubernetes", "Docker", "Terraform", "AWS", "CI/CD"],
        "preferred_skills": ["Ansible", "Prometheus", "Grafana", "Python"],
        "experience": "5+ years",
        "type": "Full-time",
    },
    {
        "title": "Frontend Lead",
        "company": "MigrantPad",
        "description": "Lead frontend development for our SaaS platform serving 1M+ users.",
        "required_skills": ["React", "TypeScript", "Next.js", "CSS", "JavaScript"],
        "preferred_skills": ["GraphQL", "Storybook", "Tailwind", "Cypress"],
        "experience": "4+ years",
        "type": "Full-time",
    },
    {
        "title": "MLOps Engineer",
        "company": "Yapı Kredi",
        "description": "Build and maintain ML infrastructure for our digital banking platform.",
        "required_skills": ["Python", "Docker", "Kubernetes", "MLflow", "AWS"],
        "preferred_skills": ["Terraform", "PyTorch", "CI/CD", "Kafka"],
        "experience": "3+ years",
        "type": "Full-time",
    },
    {
        "title": "iOS Developer",
        "company": "Banabi",
        "description": "Develop and maintain iOS application with millions of active users.",
        "required_skills": ["Swift", "UIKit", "SwiftUI", "Core Data", "Firebase"],
        "preferred_skills": ["Combine", "RxSwift", "Fastlane", "App Store Connect"],
        "experience": "3+ years",
        "type": "Full-time",
    },
    {
        "title": "Security Analyst",
        "company": "SOCRadar",
        "description": "Protect our customers from cyber threats as part of our security operations team.",
        "required_skills": ["Python", "Kubernetes", "Docker", "Penetration Testing", "SIEM"],
        "preferred_skills": ["Go", "AWS Security", "Zero Trust", "Cryptography"],
        "experience": "4+ years",
        "type": "Full-time",
    },
    {
        "title": "QA Automation Lead",
        "company": "Vodafone",
        "description": "Lead quality assurance for our digital channels serving millions of customers.",
        "required_skills": ["Selenium", "Cypress", "Python", "Java", "Docker"],
        "preferred_skills": ["Jenkins", "JMeter", "Postman", "Kubernetes"],
        "experience": "4+ years",
        "type": "Full-time",
    },
    {
        "title": "Product Manager",
        "company": "Papara",
        "description": "Drive product strategy for our fintech platform with 10M+ users.",
        "required_skills": ["Product Strategy", "User Research", "Data Analysis", "Agile", "SQL"],
        "preferred_skills": ["A/B Testing", "Figma", "JIRA", "Scrum"],
        "experience": "5+ years",
        "type": "Full-time",
    },
]

# ============================================================
# CV ÜRETİCİ
# ============================================================

def generate_cv(role):
    """Bir role şablonundan gerçekçi CV metni oluşturur."""
    skills_text = ", ".join(random.sample(role["skills"], k=random.randint(4, len(role["skills"]))))
    project = random.choice(role["projects"])

    companies = ["CompanyX", "TechCorp", "StartupY", "BigCo", "ScaleUp", "DigitalAgency", "FinTechLabs", "CloudInc"]
    company = random.choice(companies)
    prev_company = random.choice([c for c in companies if c != company])

    cv = f"""{role['title']} with {role['experience_years']} years of professional experience. Proficient in {skills_text}.

Professional Experience:
- {role['title']} at {company} ({role['experience_years']} years 2 months): {project}. Also mentored 3 junior developers and conducted technical interviews.
- Mid-level Developer at {prev_company} (2 years 6 months): Developed and maintained production systems serving 1M+ daily active users.

Education:
- {role['education']['degree']}, {role['education']['institution']}

Languages: English (Professional), Turkish (Native)"""
    return cv


def generate_parse_output(role, cv_text):
    """CV metninden parse edilmiş JSON çıktısı oluşturur."""
    skills = random.sample(role["skills"], k=random.randint(4, len(role["skills"])))
    exp_years = role["experience_years"]

    output = {
        "skills": skills,
        "experience": [
            {
                "title": role["title"],
                "company": "CompanyX",
                "start_date": f"{2026 - exp_years}-01-01",
                "end_date": "2026-07-01",
                "description": random.choice(role["projects"]),
            }
        ],
        "education": [
            {
                "degree": role["education"]["degree"],
                "field": role["education"]["degree"].split(" in ")[1] if " in " in role["education"]["degree"] else role["education"]["degree"],
                "institution": role["education"]["institution"],
                "start_year": 2026 - exp_years - 4,
                "end_year": 2026 - exp_years,
            }
        ],
        "summary": f"Experienced {role['title'].lower()} with {exp_years} years of experience. Proficient in {', '.join(skills[:4])}.",
    }
    return json.dumps(output, indent=2)


# ============================================================
# CV-JD MATCH ÜRETİCİ
# ============================================================

def generate_match_data(cv_role, jd_template):
    """Bir CV ve JD çifti için match skorları üretir."""
    cv_skills = set(cv_role["skills"])
    jd_skills = set(jd_template["required_skills"])
    jd_preferred = set(jd_template["preferred_skills"])

    matched = list(cv_skills & jd_skills)
    missing = list(jd_skills - cv_skills)
    preferred_matched = list(cv_skills & jd_preferred)

    # Skor hesapla
    if len(jd_skills) > 0:
        skill_match = len(matched) / len(jd_skills)
    else:
        skill_match = 0.5

    # Tecrübe uyumu
    cv_years = cv_role["experience_years"]
    jd_years_str = jd_template["experience"]
    jd_min_years = int(jd_years_str.split("+")[0])
    if cv_years >= jd_min_years:
        exp_score = min(1.0, cv_years / (jd_min_years + 2))
    else:
        exp_score = max(0.1, cv_years / (jd_min_years + 5))

    # Eğitim uyumu (genelde olumlu)
    edu_score = 0.7 + random.random() * 0.25

    # Genel skor
    overall = round(skill_match * 0.5 + exp_score * 0.3 + edu_score * 0.2, 2)
    overall = max(0.05, min(0.98, overall))

    # Analiz metni oluştur
    if overall >= 0.7:
        analysis = f"Strong match. The candidate has {len(matched)} of {len(jd_skills)} required skills: {', '.join(matched[:3])}. "
        if missing:
            analysis += f"Missing {', '.join(missing[:2])} but this can be learned. "
        analysis += f"With {cv_years} years of experience, the candidate exceeds the {jd_years_str} requirement. Recommended for interview."
    elif overall >= 0.35:
        analysis = f"Partial match. The candidate matches {len(matched)} of {len(jd_skills)} required skills: {', '.join(matched) if matched else 'none'}. "
        analysis += f"Missing key skills: {', '.join(missing[:3])}. Experience ({cv_years}y) is {'below' if cv_years < jd_min_years else 'meets'} the {jd_years_str} requirement. Consider if willing to upskill."
    else:
        analysis = f"Poor match. The candidate has {len(matched)} of {len(jd_skills)} required skills. "
        analysis += f"The role requires {', '.join(list(jd_skills)[:3])} but the candidate specializes in {', '.join(list(cv_skills)[:3])}. "
        analysis += f"Not recommended for this position without significant skill development."

    result = {
        "overall_score": overall,
        "skill_match_score": round(skill_match, 2),
        "experience_score": round(exp_score, 2),
        "education_score": round(edu_score, 2),
        "matched_skills": matched,
        "missing_skills": missing,
        "analysis": analysis,
    }
    return json.dumps(result, indent=2)


# ============================================================
# ANA ÜRETİCİ
# ============================================================

def generate_cv_parse_dataset(num_examples=50):
    """CV parse dataset'i oluşturur."""
    dataset = []
    instruction = "Parse the following CV text and extract structured information: skills, experience, education, and a brief summary."

    for i in range(num_examples):
        role = random.choice(ROLES)
        cv_text = generate_cv(role)
        parse_output = generate_parse_output(role, cv_text)

        dataset.append({
            "instruction": instruction,
            "input": cv_text,
            "output": parse_output,
        })

    return dataset


def generate_match_dataset(num_examples=50):
    "CV-JD match dataset'i oluşturur."
    dataset = []
    instruction = "Analyze how well this CV matches the job description. Rate each category 0.0-1.0."

    for i in range(num_examples):
        cv_role = random.choice(ROLES)
        jd = random.choice(JD_TEMPLATES)

        # Çok düşük skorlu örnekler de ekle (zıt eşleşmeler)
        if i % 5 == 4:
            zıt_cv = random.choice(ROLES)
            while zıt_cv["title"] == jd["title"]:
                zıt_cv = random.choice(ROLES)
            cv_role = zıt_cv

        cv_text = generate_cv(cv_role)
        jd_text = f"{jd['title']} at {jd['company']}\n{jd['description']}\n\nRequired: {', '.join(jd['required_skills'])}\nPreferred: {', '.join(jd['preferred_skills'])}\nExperience: {jd['experience']}\nType: {jd['type']}"

        match_output = generate_match_data(cv_role, jd)

        dataset.append({
            "instruction": instruction,
            "input": f"CV: {cv_text}\n\nJD: {jd_text}",
            "output": match_output,
        })

    return dataset


def save_dataset(dataset, filepath):
    """Dataset'i JSON dosyasına kaydeder."""
    os.makedirs(os.path.dirname(filepath), exist_ok=True)
    with open(filepath, "w", encoding="utf-8") as f:
        json.dump(dataset, f, indent=2, ensure_ascii=False)
    print(f"✅ Kaydedildi: {filepath} ({len(dataset)} örnek)")


if __name__ == "__main__":
    data_dir = os.path.join(os.path.dirname(__file__), "data")

    print("=" * 60)
    print("CV MATCHER DATASET GENERATOR")
    print("=" * 60)

    print("\n[1/2] CV Parse dataset oluşturuluyor...")
    cv_data = generate_cv_parse_dataset(50)
    save_dataset(cv_data, os.path.join(data_dir, "cv_parse_dataset.json"))

    print("\n[2/2] CV-JD Match dataset oluşturuluyor...")
    match_data = generate_match_dataset(50)
    save_dataset(match_data, os.path.join(data_dir, "cv_jd_match_dataset.json"))

    print("\n" + "=" * 60)
    print("ÖRNEK CV PARSE:")
    print("=" * 60)
    print(f"Instruction: {cv_data[0]['instruction']}")
    print(f"\nInput:\n{cv_data[0]['input'][:200]}...")
    print(f"\nOutput:\n{cv_data[0]['output'][:200]}...")

    print("\n" + "=" * 60)
    print("ÖRNEK CV-JD MATCH:")
    print("=" * 60)
    print(f"Instruction: {match_data[0]['instruction']}")
    print(f"\nInput:\n{match_data[0]['input'][:200]}...")
    print(f"\nOutput:\n{match_data[0]['output'][:200]}...")

    print("\n✅ Dataset hazır! Fine-tuning'e geçebilirsin.")
