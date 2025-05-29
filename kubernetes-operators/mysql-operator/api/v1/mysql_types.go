package v1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MySQLSpec struct {
    Image        string `json:"image"`
    Database     string `json:"database"`
    Password     string `json:"password"`
    StorageSize  string `json:"storage_size"`
}

type MySQLStatus struct {}

type MySQL struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   MySQLSpec   `json:"spec,omitempty"`
    Status MySQLStatus `json:"status,omitempty"`
}

type MySQLList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []MySQL `json:"items"`
}
