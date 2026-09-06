UPDATE password_resets SET used_at = now(), status = 'used', updated_at = now()
WHERE id = $1
