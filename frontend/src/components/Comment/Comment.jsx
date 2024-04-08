// import userImage from "/src/static/img/user_image.png";
import "./Comment.css"
const Comment= ({comment, onDelete}) => {

    const handleDeleteCommentClick = () => {
        onDelete(comment._id);
    };

    return (
        <div className="comment-info">
            <div className="comment-user">
                <p>{comment.User.username}</p>
                <img className="user-image" src={comment.User.image} alt="image" />
            </div>
            <div className="create-comment">
                <p>{comment.message}</p>
                <button onClick={handleDeleteCommentClick}>Delete</button>
            </div>
        </div>
    );
};

export default Comment;