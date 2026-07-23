package llm

const CVParsePrompt = `You are a CV/resume parsing expert. Extract structured information from the CV text below.

Return a JSON object with these fields:
- skills: array of strings listing all technical and soft skills found
- experience: array of objects with { title, company, start_date, end_date, description }
- education: array of objects with { degree, field, institution, start_year, end_year }
- summary: a 2-3 sentence professional summary of the candidate

CV TEXT:
%s

Return ONLY valid JSON, no other text.`

const JDAnalyzePrompt = `You are a job description analyst. Extract key information from the job posting below.

Return a JSON object with these fields:
- required_skills: array of strings, skills that are explicitly required
- preferred_skills: array of strings, skills that are preferred/optional
- experience_level: the required experience level (Junior, Mid, Senior, Lead)
- employment_type: Full-time, Part-time, Contract, Remote
- summary: 2-3 sentence summary of the role

JOB DESCRIPTION:
%s

Return ONLY valid JSON, no other text.`

const CVJDMatchPrompt = `You are a hiring expert analyzing how well a candidate's CV matches a job description.

CV SUMMARY:
- Skills: %s
- Experience: %s
- Education: %s
- Summary: %s

JOB REQUIREMENTS:
- Required Skills: %s
- Preferred Skills: %s
- Experience Level: %s
- Description: %s

Analyze the match and return a JSON object with:
- overall_score: 0-100 overall match percentage
- skill_match_score: 0-100 how well skills match
- experience_score: 0-100 how well experience matches
- education_score: 0-100 how well education matches
- matched_skills: array of skills the candidate has that match
- missing_skills: array of required skills the candidate lacks
- analysis: 3-4 sentence detailed analysis of the match, including strengths and gaps

Return ONLY valid JSON, no other text.`
