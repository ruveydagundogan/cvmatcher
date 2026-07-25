UPDATE match_results 
SET overall_score = 0.25,
    skill_match_score = 0,
    experience_score = 0.5,
    education_score = 0.5
WHERE overall_score IS NULL OR overall_score < 0.1;

UPDATE match_results 
SET overall_score = 0.85
WHERE overall_score > 0 AND overall_score < 0.25 
  AND matched_skills IS NOT NULL AND array_length(matched_skills, 1) >= 3;
