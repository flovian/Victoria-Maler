package models

type Evidence struct {
    ID        int64  `db:"id" json:"id"`
    ReportID  int64  `db:"report_id" json:"report_id"`
    FilePath  string `db:"file_path" json:"file_path"`
    Timestamp int64  `db:"timestamp" json:"timestamp"`
}
